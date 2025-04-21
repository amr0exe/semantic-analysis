package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GenerateQuestions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var reqBody GenRequest
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Fatal(err)
		}

		ctx, collection, cancel := InitDB()
		defer cancel()

		subID := reqBody.SubjectID
		subjectID, err := primitive.ObjectIDFromHex(subID) // chemistry
		if err != nil {
			log.Fatal("Invalid subjectID", err)
		}

		var gradeID *primitive.ObjectID
		grdID := reqBody.GradeID
		if grdID != "" {
			parsedGradeID, err := primitive.ObjectIDFromHex(grdID)
			if err != nil {
				log.Fatal("Invalid gradeID", err)
			}
			gradeID = &parsedGradeID
		}

		qAno := reqBody.GroupA
		qBno := reqBody.GroupB
		qCno := reqBody.GroupC

		//finaldata := make(map[string][]SingleQuestion)
		var finaldata Prequest

		finaldata.SubjectID = subjectID
		if gradeID != nil {
			finaldata.GradeID = *gradeID
		}

		groupAQuestions, err := fetchRandomQuestions(collection, "GroupA", qAno, subjectID, gradeID, ctx)
		if err != nil {
			log.Fatal(err)
		}
		//finaldata["GroupA"] = groupAQuestions
		finaldata.QuestionsBlock.GroupA = groupAQuestions

		groupBQuestions, err := fetchRandomQuestions(collection, "GroupB", qBno, subjectID, gradeID, ctx)
		if err != nil {
			log.Fatal(err)
		}
		//finaldata["GroupB"] = groupBQuestions
		finaldata.QuestionsBlock.GroupB = groupBQuestions

		groupCQuestions, err := fetchRandomQuestions(collection, "GroupC", qCno, subjectID, gradeID, ctx)
		if err != nil {
			log.Fatal(err)
		}
		//finaldata["GroupC"] = groupCQuestions
		finaldata.QuestionsBlock.GroupC = groupCQuestions

		// ----

		// check-for-duplicates
		checkResp, err := checkDuplicates(finaldata)
		if err != nil {
			log.Fatal("Error checking duplicates", err)
		}

		if checkResp.Status == "ok" {
			saveToFile(finaldata)
			fmt.Printf("No duplicates, File Written Successfully")

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"data":    finaldata,
			})
		}

		var updates []Opayload
		if checkResp.Status == "conflict" {
			log.Printf("Duplictes occured.. processing...")

			for _, dup := range checkResp.Duplicates {
				groupQuestions, err := fetchRandomQuestions(collection, dup.Group, 1, subjectID, gradeID, ctx)
				if err != nil {
					log.Fatal("Error fetching questions for updates", err)
				}

				// local-update
				switch dup.Group {
				case "GroupA":
					finaldata.QuestionsBlock.GroupA[dup.Index] = groupQuestions[0]
				case "GroupB":
					finaldata.QuestionsBlock.GroupB[dup.Index] = groupQuestions[0]
				case "GroupC":
					finaldata.QuestionsBlock.GroupC[dup.Index] = groupQuestions[0]
				}

				updates = append(updates, Opayload{
					Group:       dup.Group,
					Index:       dup.Index,
					NewQuestion: groupQuestions[0],
				})
			}

			err = updateDuplicate(OutPayload{Updates: updates})
			if err != nil {
				log.Printf("Error updating duplicates")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]any{
					"success": false,
					"message": "failed updating Data",
					"error":   err.Error(),
				})

				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]any{
				"success": true,
				"message": "updated Data",
				"data":    finaldata,
			})
		}

	}
}

func fetchRandomQuestions(collection *mongo.Collection, groupField string, sampleSize int, subjectID primitive.ObjectID, gradeID *primitive.ObjectID, ctx context.Context) ([]SingleQuestion, error) {
	filter := bson.M{
		"subjectID": subjectID,
	}
	if gradeID != nil {
		filter["gradeID"] = *gradeID
	}

	pipeline := mongo.Pipeline{
		{{
			Key: "$match", Value: filter,
		}},
		{{
			Key: "$project", Value: bson.M{
				"groupQuestions": fmt.Sprintf("$questionsBlock.%s", groupField),
			},
		}},

		{{
			Key: "$unwind", Value: "$groupQuestions",
		}},

		{{
			Key: "$sample", Value: bson.M{
				"size": sampleSize,
			},
		}},

		{{
			Key: "$project", Value: bson.M{
				"associatedMarks": "$groupQuestions.associatedMarks",
				"editorContent":   "$groupQuestions.editorContent",
				"questionType":    "$groupQuestions.questionType",
				"title":           "$groupQuestions.title",
			},
		}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var question []SingleQuestion
	if err := cursor.All(ctx, &question); err != nil {
		return nil, err
	}

	//for i, q := range question {
	//	question[i].Title = extractEnglishText(q.Title)
	//}

	return question, nil
}

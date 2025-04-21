package main

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prequest struct {
	SubjectID		primitive.ObjectID		`bson:"subjectID" json:"subjectID"`
	GradeID			primitive.ObjectID		`bson:"gradeID" json:"gradeID"`
	QuestionsBlock	QuestsBlocks			`bson:"questionsBlock" json:"questionsBlock"`
}

type QuestsBlocks struct {
	GroupA []SingleQuestion 	`bson:"GroupA" json:"GroupA"`
	GroupB []SingleQuestion 	`bson:"GroupB" json:"GroupB"`
	GroupC []SingleQuestion 	`bson:"GroupC" json:"GroupC"`
 }

// for incoming /api/questions-generate
type GenRequest struct {
	SubjectID	string		`json:"subjectID"`
	GradeID		string		`json:"gradeID"`
	GroupA		int			`json:"groupA"`
	GroupB		int			`json:"groupB"`
	GroupC		int			`json:"groupC"`
}

// for editorContents section
type SingleQuestion struct {
	AssociatedMarks	string			`bson:"associatedMarks" json:"associatedMarks"`
	EditorContent 	EditorJsContent	`bson:"editorContent" json:"editorContent"`
	QuestionType 	string			`bson:"questionType" json:"questionType"`
	Title 			string 			`bson:"title" json:"title"`
}

type EditorJsContent struct {
	Blocks	[]Block	`bson:"blocks" json:"blocks"`
	Time 	int64	`bson:"time" json:"time"`
	Version	string	`bson:"version" json:"version"`
}

type Block struct {
	ID 		string	`json:"id"`	
	Type	string	`json:"type"`
	Data	DataB	`bson:"data" json:"data"`
}

type DataB struct {
	Text	string 	`bson:"text" json:"text"`
}

// for outbound-req
type CheckDuplicateResponse struct {
	Duplicates	[]Duplicate		`json:"duplicates"`
	Message		string			`json:"message"`
	Status		string			`json:"status"`
}
type Duplicate struct {
	Group     string `json:"group"`
	Index     int    `json:"index"`
	SimilarTo struct {
		Group           string
		Index           int
		SimilarityScore float64
		Title           string
	} `json:"similar_to"`
	Title string `json:"title"`
}

// for outbound /update
type OutPayload struct {
	Updates []Opayload	`json:"updates"`
}

type Opayload struct {
	Group 		string	`json:"group"`
	Index		int		`json:"index"`
	NewQuestion	SingleQuestion	`json:"new_question"`
}

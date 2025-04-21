import numpy as np
import re
from flask import Flask, request, jsonify
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

# Initialize Flask app
app = Flask(__name__)

# Load the pre-trained SentenceTransformer model
model = SentenceTransformer("all-mpnet-base-v2")

# In-memory storage (you could replace this with a database later)
stored_groups = {}


# Preprocessing and embedding generation function
def preprocess_and_embed(sentences):
    preprocessed_sentences = [
        re.sub(r'^(State|Define|Explain|How|What is|Describe|Identify)\s+', '', sentence).lower()
        for sentence in sentences
    ]
    embeddings = model.encode(preprocessed_sentences)
    return embeddings


def find_similar_questions(groups, extremely_similar_threshold=0.90, quite_similar_threshold=0.70):
    all_questions = []
    group_names = list(groups.keys())

    # Flatten all questions into a list with group and index info
    for group in group_names:
        for idx, question in enumerate(groups[group]):
            all_questions.append((group, idx, question['title']))

    sentences = [q[2] for q in all_questions]  # Only titles
    embeddings = preprocess_and_embed(sentences)

    similarities = cosine_similarity(embeddings)
    np.fill_diagonal(similarities, -1)

    duplicates = []
    seen_pairs = set()  # <-- NEW: to track processed pairs

    for idx, (group_a, index_a, question_a) in enumerate(all_questions):
        similar_idx = np.argmax(similarities[idx])
        similarity_score = similarities[idx][similar_idx]

        if similarity_score >= quite_similar_threshold:
            group_b, index_b, question_b = all_questions[similar_idx]

            # Create a consistent pair (smaller index first)
            pair = tuple(sorted([(group_a, index_a), (group_b, index_b)]))

            # If already seen, skip
            if pair in seen_pairs:
                continue

            # Otherwise, mark as seen and add to duplicates
            seen_pairs.add(pair)

            duplicates.append({
                'group': group_a,
                'index': index_a,
                'title': question_a,
                'similar_to': {
                    'group': group_b,
                    'index': index_b,
                    'title': question_b,
                    'similarity_score': float(similarity_score)
                }
            })

    return duplicates


# POST endpoint: Check for duplicate questions
@app.route('/check-duplicates', methods=['POST'])
def check_duplicates():
    try:
        data = request.get_json()

        if not data:
            return jsonify({'status': 'error', 'message': 'No data provided'}), 400

        questions_block = data.get('questionsBlock')
        if not questions_block:
            return jsonify({ 'status': 'false', 'message': 'questionsblock not found'})

        global stored_groups
        stored_groups = questions_block

        duplicates = find_similar_questions(stored_groups)

        if not duplicates:
            return jsonify({
                'status': 'ok',
                'message': 'All questions are unique!'
            }), 200
        else:
            return jsonify({
                'status': 'conflict',
                'message': 'Similar questions found!',
                'duplicates': duplicates
            }), 200

    except Exception as e:
        return jsonify({'status': 'error', 'message': str(e)}), 500


# POST endpoint: Update duplicate questions
@app.route('/update-duplicates', methods=['POST'])
def update_duplicates():
    try:
        updates = request.get_json()

        if not updates or 'updates' not in updates:
            return jsonify({'status': 'error', 'message': 'No updates provided'}), 400

        global stored_groups

        for update in updates['updates']:
            group = update['group']
            index = update['index']
            new_question = update['new_question']

            if group in stored_groups and index < len(stored_groups[group]):
                stored_groups[group][index] = new_question
            else:
                return jsonify({'status': 'error', 'message': f'Invalid group or index: {group}, {index}'}), 400

        # Optionally recheck after update
        duplicates = find_similar_questions(stored_groups)

        if not duplicates:
            return jsonify({
                'status': 'ok',
                'message': 'Updates successful. All questions are now unique!',
                'data': stored_groups
            }), 200
        else:
            return jsonify({
                'status': 'conflict',
                'message': 'Updates applied but still found some similar questions.',
                'duplicates': duplicates
            }), 200

    except Exception as e:
        return jsonify({'status': 'error', 'message': str(e)}), 500


# Run the Flask app
if __name__ == '__main__':
    app.run(debug=True)

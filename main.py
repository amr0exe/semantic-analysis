import numpy as np
import re
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity

model = SentenceTransformer("all-mpnet-base-v2")

sentences = [
    "Explain how photosynthesis enables plants to store solar energy as chemical energy.",
    "Describe the process through which plants capture and convert sunlight into stored energy.",
    "What is the role of chlorophyll in the photosynthesis process?",
    "Define Newton’s second law of motion and explain its relationship with force and acceleration.",
    "What is Newton's second law, and how does it relate to mass and force?",
    "Describe how gravitational potential energy changes with height.",
    "What are the major organelles found in a plant cell and their functions?",
    "Identify the key structures inside a plant cell and describe their roles.",
    "How do different wavelengths of light affect the rate of photosynthesis?",
    "What is the difference between an ionic bond and a covalent bond?",
    "Explain the steps involved in the nitrogen cycle.",
    "Describe how nitrogen moves through the environment during the nitrogen cycle.",
    "How does air pressure influence weather patterns?",
    "What is the difference between renewable and nonrenewable resources?",
    "Define electric resistance and its role in a circuit.",
    "What is electrical resistance and how does it affect current flow?",
    "How do vaccines create immunity in the human body?",
    "What causes earthquakes and how are they measured?",
    "Explain how fossils provide evidence for evolution.",
    "What is osmosis and why is it vital for cell survival?",
    "Describe the process of osmosis and its importance to living cells.",
    "How do chemical catalysts speed up reactions without being consumed?",
    "What factors influence the rate of evaporation of water?",
    "Define inertia and how it affects moving objects.",
    "What is inertia and how does it influence the motion of objects?",
    "How does sound travel differently through solids, liquids, and gases?",
    "What are the characteristics that distinguish prokaryotic cells from eukaryotic cells?",
    "Describe how antibiotics target bacterial cells but not human cells.",
    "What is the principle behind hydraulic systems?",
    "How does genetic variation contribute to the process of natural selection?",
    "What is the role of decomposers in an ecosystem?",
    "Explain how the law of conservation of energy applies to a roller coaster ride.",
    "How are metamorphic rocks formed from existing rock types?",
    "What is the significance of the carbon footprint in environmental science?",
    "Describe how the periodic table is organized based on atomic structure.",
    "What is the Doppler effect and how does it apply to sound and light waves?",
    "How does the nervous system transmit messages throughout the body?",
    "What are the major phases of mitosis and what occurs in each phase?",
    "Explain the relationship between wavelength and frequency in the electromagnetic spectrum.",
    "What happens during the process of nuclear fusion inside the Sun?",
    "How does the structure of a virus differ from that of a bacterial cell?"
]



# Or try preprocessing before encoding
preprocessed_sentences = [
    re.sub(r'^(State|Define|Explain|How|What is|Describe|Identify)\s+', '', sentence).lower()
    for sentence in sentences
]

embeddings = model.encode(preprocessed_sentences)

similarities = cosine_similarity(embeddings)

np.fill_diagonal(similarities, -1)


extremely_similar_threshold = 0.90
quite_similar_threshold = 0.80

similarity_data = []
for idx, row in enumerate(similarities):
    most_similar_idx = np.argmax(row)
    similarity_score = row[most_similar_idx]

    if similarity_score > 0:
        similarity_data.append((idx, most_similar_idx, similarity_score))

similarity_data.sort(key=lambda x:x[2], reverse=True)

with open("output.txt", "w", encoding="utf-8") as f:
    for idx, most_similar_idx, similarity_score in similarity_data:
        f.write(f"\n🔹 Sentence: {sentences[idx]}\n")
        f.write(f"+ Most similar to: {sentences[most_similar_idx]}\n")
        f.write(f"Similarity Score: {similarity_score:.4f}\n")

        # Find extremely similar ones
        row = similarities[idx]
        extremely_similar = np.where(row >= extremely_similar_threshold)[0]
        for sim_idx in extremely_similar:
            if sim_idx != most_similar_idx:
                f.write(f"- Extremely Similar: {sentences[sim_idx]} (Score: {row[sim_idx]:.4f})\n")

        # Find quite similar ones
        quite_similar = np.where((row >= quite_similar_threshold) & (row < extremely_similar_threshold))[0]
        for sim_idx in quite_similar:
            f.write(f"  Quite Similar: {sentences[sim_idx]} (Score: {row[sim_idx]:.4f})\n")
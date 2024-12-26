from transformers import pipeline

# Load text classification pipeline
classification_pipeline = pipeline(
    "text-classification", model="distilbert-base-uncased"
)


def categorize_text(prompt: str) -> str:
    """
    Categorize the given text using a mini LLM.
    """
    result = classification_pipeline(prompt)
    return result[0]["label"]

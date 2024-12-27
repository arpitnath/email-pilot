from app.services.model_loader import get_classification_pipeline

# Initialize classification pipeline
classification_pipeline = get_classification_pipeline()


def categorize_text(prompt: str) -> str:
    """
    Categorize the given text using the classification pipeline.
    """
    result = classification_pipeline(prompt)
    return result[0]["label"]

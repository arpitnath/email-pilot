from app.services.model_loader import get_sentiment_pipeline

# Initialize sentiment pipeline
sentiment_pipeline = get_sentiment_pipeline()


def analyze_sentiment(prompt: str) -> str:
    """
    Analyze the sentiment of the given text using the sentiment pipeline.
    """
    result = sentiment_pipeline(prompt)
    return result[0]["label"]

from transformers import pipeline

# Load sentiment analysis pipeline
sentiment_pipeline = pipeline("sentiment-analysis")


def analyze_sentiment(prompt: str) -> str:
    """
    Analyze the sentiment of the given text using a mini LLM.
    """
    result = sentiment_pipeline(prompt)
    return result[0]["label"]

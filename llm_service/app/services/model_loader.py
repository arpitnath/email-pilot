import os
from transformers import AutoTokenizer, AutoModelForSeq2SeqLM, pipeline

# Load model names from environment variables
LARGE_MODEL = os.getenv("LARGE_MODEL", "facebook/bart-large-cnn")
MINI_MODEL_1 = os.getenv("MINI_MODEL_1", "distilbert-base-uncased")
MINI_MODEL_2 = os.getenv(
    "MINI_MODEL_2", "nlptown/bert-base-multilingual-uncased-sentiment"
)


def load_pipeline(task: str, model_name: str):
    """
    Loads a pipeline for the given task and model name.
    """
    return pipeline(task, model=model_name)


def get_summarization_pipeline():
    """
    Returns the pipeline for summarization tasks.
    """
    return load_pipeline("summarization", LARGE_MODEL)


def get_classification_pipeline():
    """
    Returns the pipeline for text classification tasks.
    """
    return load_pipeline("text-classification", MINI_MODEL_1)


def get_sentiment_pipeline():
    """
    Returns the pipeline for sentiment analysis tasks.
    """
    return load_pipeline("sentiment-analysis", MINI_MODEL_2)

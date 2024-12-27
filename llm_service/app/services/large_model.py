from app.services.model_loader import get_summarization_pipeline

# Initialize summarization pipeline
summarization_pipeline = get_summarization_pipeline()


def summarize_text(prompt: str) -> str:
    """
    Summarize the given text using the summarization pipeline.
    """
    result = summarization_pipeline(
        prompt, max_length=50, min_length=25, do_sample=False
    )
    return result[0]["summary_text"]

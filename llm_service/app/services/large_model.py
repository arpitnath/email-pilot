from app.services.model_loader import get_summarization_pipeline

# Initialize summarization pipeline
summarization_pipeline = get_summarization_pipeline()


def summarize_text(prompt: str) -> str:
    """
    Summarize the given text using the summarization pipeline.
    """
    input_length = len(prompt.split())
    max_length = max(10, input_length * 3 // 2)

    if input_length > 1000:
        return "Input too long"

    result = summarization_pipeline(
        prompt, max_length=max_length, min_length=5, truncation=True
    )
    return result[0]["summary_text"]

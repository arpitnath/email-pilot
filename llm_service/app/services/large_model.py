from transformers import pipeline

# Load summarization pipeline
summarization_pipeline = pipeline("summarization", model="EleutherAI/gpt-neo-1.3B")


def summarize_text(prompt: str) -> str:
    """
    Summarize the given text using a large LLM.
    """
    result = summarization_pipeline(
        prompt, max_length=50, min_length=25, do_sample=False
    )
    return result[0]["summary_text"]

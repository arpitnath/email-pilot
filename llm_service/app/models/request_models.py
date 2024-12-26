from pydantic import BaseModel


class SummarizationRequest(BaseModel):
    prompt: str  # The text to summarize


class CategorizationRequest(BaseModel):
    prompt: str  # The text to categorize

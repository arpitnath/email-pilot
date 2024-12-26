from pydantic import BaseModel


class SummarizationResponse(BaseModel):
    summary: str  # The summarized text


class CategorizationResponse(BaseModel):
    category: str  # The category of the text

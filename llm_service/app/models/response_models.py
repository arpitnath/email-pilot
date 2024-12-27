from pydantic import BaseModel


class SummarizationResponse(BaseModel):
    summary: str


class CategorizationResponse(BaseModel):
    category: str


class SentimentAnalysisResponse(BaseModel):
    sentiment: str

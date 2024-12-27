from pydantic import BaseModel


class SummarizationRequest(BaseModel):
    prompt: str


class CategorizationRequest(BaseModel):
    prompt: str


class SentimentAnalysisRequest(BaseModel):
    prompt: str

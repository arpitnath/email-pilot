from fastapi import APIRouter, HTTPException
from app.services.mini_model_2 import analyze_sentiment
from app.models.request_models import SentimentAnalysisRequest
from app.models.response_models import SentimentAnalysisResponse

router = APIRouter()


@router.post("/sentiment", response_model=SentimentAnalysisResponse)
async def sentiment_analysis(request: SentimentAnalysisRequest):
    """
    Analyze the sentiment of the given text prompt.
    """
    try:
        result = analyze_sentiment(request.prompt)
        return {"sentiment": result}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

from fastapi import APIRouter, HTTPException
from app.models.request_models import SummarizationRequest
from app.models.response_models import SummarizationResponse
from app.services.large_model import summarize_text

router = APIRouter()


@router.post("/summarize", response_model=SummarizationResponse)
def summarize(request: SummarizationRequest):
    try:
        summary = summarize_text(request.prompt)
        return {"summary": summary}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error summarizing text: {str(e)}")

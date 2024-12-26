from fastapi import APIRouter, HTTPException
from app.models.request_models import CategorizationRequest
from app.models.response_models import CategorizationResponse
from app.services.mini_model_1 import categorize_text

router = APIRouter()


@router.post("/categorize", response_model=CategorizationResponse)
def categorize(request: CategorizationRequest):
    try:
        category = categorize_text(request.prompt)
        return {"category": category}
    except Exception as e:
        raise HTTPException(
            status_code=500, detail=f"Error categorizing text: {str(e)}"
        )

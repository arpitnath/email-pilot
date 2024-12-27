from fastapi import FastAPI
from app.routers import summarize, categorize, sentiment_analyzer


app = FastAPI(
    title="LLM Service",
    description="Service for Summarization and Categorization using LLMs",
    version="1.0.0",
)


app.include_router(summarize.router, prefix="/api")
app.include_router(categorize.router, prefix="/api")
app.include_router(sentiment_analyzer.router, prefix="/api")


# Health check route
@app.get("/")
def health_check():
    return {"message": "LLM Service is running"}


if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)

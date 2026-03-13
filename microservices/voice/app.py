from fastapi import FastAPI, UploadFile, File, HTTPException
import uvicorn
import whisper
import io

app = FastAPI(title="Voice Service")

# Загрузка модели Whisper
model = whisper.load_model("base")

@app.get("/")
async def root():
    return {"message": "Voice Service is running"}

@app.get("/health")
async def health():
    return {"status": "ok"}

@app.post("/transcribe")
async def transcribe_audio(file: UploadFile = File(...)):
    try:
        # Читаем аудиофайл
        audio_bytes = await file.read()
        
        # Сохраняем во временный файл
        with open("temp_audio.wav", "wb") as f:
            f.write(audio_bytes)
        
        # Распознаем речь
        result = model.transcribe("temp_audio.wav", language="ru")
        
        return {
            "text": result["text"],
            "language": result["language"],
            "segments": result["segments"]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@app.post("/process-prescription")
async def process_prescription(file: UploadFile = File(...)):
    try:
        audio_bytes = await file.read()
        
        with open("temp_audio.wav", "wb") as f:
            f.write(audio_bytes)
        
        result = model.transcribe("temp_audio.wav", language="ru")
        
        # Здесь можно добавить парсинг текста для извлечения информации о лекарствах
        text = result["text"].lower()
        
        # Простой парсинг (для демо)
        medications = []
        if "амоксициллин" in text:
            medications.append({
                "name": "Амоксициллин",
                "dosage": "500 мг",
                "frequency": "3 раза в день",
                "duration": "7 дней"
            })
        if "парацетамол" in text:
            medications.append({
                "name": "Парацетамол",
                "dosage": "500 мг",
                "frequency": "при боли",
                "duration": "по необходимости"
            })
        
        return {
            "text": text,
            "medications": medications,
            "confidence": 0.95
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)

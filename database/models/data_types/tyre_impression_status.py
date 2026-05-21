import enum


class TyreImpressionStatus(enum.Enum):
    uploaded = "uploaded"
    processing = "processing"
    processed = "processed"
    matched = "matched"
    failed = "failed"

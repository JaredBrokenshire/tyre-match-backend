import enum


class TyreImpressionStatus(enum.Enum):
    uploaded = "uploaded"
    queued = "queued"
    preprocessing = "preprocessing"
    preprocessed = "preprocessed"
    extracting_features = "extracting_features"
    matching= "matching"
    matched = "matched"
    failed = "failed"

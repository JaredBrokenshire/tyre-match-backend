import enum


class TyreImpressionStatus(enum.Enum):
    uploaded = "uploaded"
    processing = "processing"
    processed = "processed"
    matched = "matched"
    failed = "failed"

    @classmethod
    def terminal_states(cls):
        return {cls.processed, cls.matched, cls.failed}
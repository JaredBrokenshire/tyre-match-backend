def slim_tyre_model_response(tyre_model):
    return {
        "id": tyre_model.id,
        "manufacturer": tyre_model.manufacturer,
        "model_name": tyre_model.model_name,
        "category": tyre_model.category,
        "vehicle_type": tyre_model.vehicle_type,
    }

def tyre_model_response(tyre_model):
    return {
        "id": tyre_model.id,
        "manufacturer": tyre_model.manufacturer,
        "model_name": tyre_model.model_name,
        "category": tyre_model.category,
        "vehicle_type": tyre_model.vehicle_type,
        "width_mm": tyre_model.width_mm,
        "aspect_ratio": tyre_model.aspect_ratio,
        "rim_diameter_inches": tyre_model.rim_diameter_inches,
        "groove_count": tyre_model.groove_count,
        "pattern_type": tyre_model.pattern_type,
        "tread_pitch_length_mm": tyre_model.tread_pitch_length_mm,
        "dataset_source": tyre_model.dataset_source,
        "notes": tyre_model.notes,
    }
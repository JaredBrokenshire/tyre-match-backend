def paginated_response(data, total_count: int):
    return {
        "data": data,
        "total_count": total_count,
    }

def error_response(status, error):
    return {
        "error": error,
    }, status
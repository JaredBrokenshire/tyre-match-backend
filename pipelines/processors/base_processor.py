import uuid
import logging
from database.models.data_types.files import FileModel
from domain.exceptions import FileSaveError, DatabaseError
from services.file_service import ProcessedImageRequest, FileService

logger = logging.getLogger(__name__)

class BaseProcessor:
    name: str

    def __init__(self):
        self.file_service = FileService()


    def process(self, input_file_path: str, context: dict) -> str:
        """
        :param input_file_path: Input file path
        :param context: Context parameters
        :return: Output file path
        """
        pass

    def transform(self, image, context: dict):
        """
        :param image: Input image
        :param context: Context parameters
        :return: Output image
        """
        pass

    def save_file(self, image, context: dict, stage_name: str) -> str:
        request = ProcessedImageRequest(
            image=image,
            file_name=f"{uuid.uuid4()}.jpg",
            upload_directory=context["output_directories"][stage_name],
            model=FileModel.tyre_impression,
            model_id=context["processing_id"],
            file_type=context["file_types_on_completion"][stage_name],
        )

        try:
            file_record = self.file_service.save_processed_image(request)
        except PermissionError as e:
            logger.error(f"PermissionError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"PermissionError from file service in {stage_name} processor: {e}")
        except OSError as e:
            logger.error(f"OSError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"OSError from file service in {stage_name} processor: {e}")
        except FileSaveError as e:
            logger.error(f"FileSaveError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"FileSaveError from file service in {stage_name} processor: {e}")
        except DatabaseError as e:
            logger.error(f"DatabaseError from file service in {stage_name} processor: {e}")
            raise FileSaveError(f"DatabaseError from file service in {stage_name} processor: {e}")

        return file_record.file_location
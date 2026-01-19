import math
import os


def get_rounding_increment() -> int | None:
    """Get the duration rounding increment from environment variable.

    Returns:
        None if no rounding, or 15/30 for the increment in minutes.

    Raises:
        ValueError: If TACT_DURATION_ROUNDING has an invalid value.
    """
    value = os.getenv("TACT_DURATION_ROUNDING", "none").lower()

    if value == "none":
        return None
    elif value == "15":
        return 15
    elif value == "30":
        return 30
    else:
        raise ValueError(
            f"Invalid TACT_DURATION_ROUNDING value: {value}. "
            "Must be 'none', '15', or '30'."
        )


def round_duration(minutes: int | None, increment: int | None = None) -> int | None:
    """Round duration up to the nearest increment.

    Args:
        minutes: The duration in minutes, or None if not set.
        increment: The rounding increment (15 or 30), or None for no rounding.
                   If not provided, reads from TACT_DURATION_ROUNDING env var.

    Returns:
        The rounded duration, or None if input was None.
    """
    if minutes is None:
        return None

    if increment is None:
        increment = get_rounding_increment()

    if increment is None:
        return minutes

    # Round up to nearest increment using ceiling division
    return math.ceil(minutes / increment) * increment

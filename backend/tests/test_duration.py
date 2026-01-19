import os

import pytest

from tact.utils.duration import get_rounding_increment, round_duration


class TestGetRoundingIncrement:
    def test_default_is_none(self, monkeypatch):
        monkeypatch.delenv("TACT_DURATION_ROUNDING", raising=False)
        assert get_rounding_increment() is None

    def test_none_value(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "none")
        assert get_rounding_increment() is None

    def test_none_case_insensitive(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "NONE")
        assert get_rounding_increment() is None

    def test_15_minutes(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "15")
        assert get_rounding_increment() == 15

    def test_30_minutes(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "30")
        assert get_rounding_increment() == 30

    def test_invalid_value_raises(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "10")
        with pytest.raises(ValueError) as exc_info:
            get_rounding_increment()
        assert "Invalid TACT_DURATION_ROUNDING value" in str(exc_info.value)

    def test_invalid_string_raises(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "always")
        with pytest.raises(ValueError):
            get_rounding_increment()


class TestRoundDuration:
    def test_none_input_returns_none(self):
        assert round_duration(None, 15) is None
        assert round_duration(None, 30) is None
        assert round_duration(None, None) is None

    def test_no_rounding_when_increment_none(self):
        assert round_duration(7, None) == 7
        assert round_duration(16, None) == 16
        assert round_duration(45, None) == 45

    # 15-minute rounding tests
    def test_15min_zero(self):
        assert round_duration(0, 15) == 0

    def test_15min_round_up_from_1(self):
        assert round_duration(1, 15) == 15

    def test_15min_round_up_from_7(self):
        assert round_duration(7, 15) == 15

    def test_15min_round_up_from_14(self):
        assert round_duration(14, 15) == 15

    def test_15min_exact_boundary(self):
        assert round_duration(15, 15) == 15

    def test_15min_round_up_from_16(self):
        assert round_duration(16, 15) == 30

    def test_15min_round_up_from_29(self):
        assert round_duration(29, 15) == 30

    def test_15min_exact_30(self):
        assert round_duration(30, 15) == 30

    def test_15min_round_up_from_31(self):
        assert round_duration(31, 15) == 45

    def test_15min_exact_45(self):
        assert round_duration(45, 15) == 45

    def test_15min_round_up_from_46(self):
        assert round_duration(46, 15) == 60

    def test_15min_exact_60(self):
        assert round_duration(60, 15) == 60

    # 30-minute rounding tests
    def test_30min_zero(self):
        assert round_duration(0, 30) == 0

    def test_30min_round_up_from_1(self):
        assert round_duration(1, 30) == 30

    def test_30min_round_up_from_7(self):
        assert round_duration(7, 30) == 30

    def test_30min_round_up_from_29(self):
        assert round_duration(29, 30) == 30

    def test_30min_exact_boundary(self):
        assert round_duration(30, 30) == 30

    def test_30min_round_up_from_31(self):
        assert round_duration(31, 30) == 60

    def test_30min_round_up_from_45(self):
        assert round_duration(45, 30) == 60

    def test_30min_exact_60(self):
        assert round_duration(60, 30) == 60

    def test_30min_round_up_from_61(self):
        assert round_duration(61, 30) == 90

    # Test with env var (no explicit increment)
    def test_reads_env_var_when_increment_not_provided(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "15")
        assert round_duration(7) == 15

    def test_no_rounding_from_env_var(self, monkeypatch):
        monkeypatch.setenv("TACT_DURATION_ROUNDING", "none")
        assert round_duration(7) == 7

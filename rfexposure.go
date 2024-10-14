package rfexposure

import (
	"math"
)

func TestStub() []float64 {

	var xmtr_power int16 = 1000
	var feedline_length int16 = 73
	var duty_cycle float64 = 0.5
	var per_30 float64 = 0.5

	c1 := CableValues{
		k1: 0.122290,
		k2: 0.000260,
	}

	all_frequency_values := []FrequencyValues{
		{
			freq:    7.3,
			swr:     2.25,
			gaindbi: 1.5,
		},
		{
			freq:    14.35,
			swr:     1.35,
			gaindbi: 1.5,
		},
		{
			freq:    18.1,
			swr:     3.7,
			gaindbi: 1.5,
		},
		{
			freq:    21.45,
			swr:     4.45,
			gaindbi: 1.5,
		},
		{
			freq:    24.99,
			swr:     4.1,
			gaindbi: 1.5,
		},
		{
			freq:    29.7,
			swr:     2.18,
			gaindbi: 4.5,
		},
	}

	var uncontrolled_safe_distances []float64

	for _, f := range all_frequency_values {
		uncontrolled_safe_distances = append(uncontrolled_safe_distances, CalculateUncontrolledSafeDistance(f, c1, xmtr_power, feedline_length, duty_cycle, per_30))
	}

	return uncontrolled_safe_distances
}

type CableValues struct {
	k1 float64
	k2 float64
}

type FrequencyValues struct {
	freq    float64
	swr     float64
	gaindbi float64
}

func CalculateUncontrolledSafeDistance(freq_values FrequencyValues, cable_values CableValues, transmitter_power int16,
	feedline_length int16, duty_cycle float64, uncontrolled_percentage_30_minutes float64) float64 {

	gamma := CalculateReflectionCoefficient(freq_values)

	feedline_loss_per_100ft_at_frequency := CalculateFeedlineLossPer100ftAtFrequency(freq_values, cable_values)

	feedline_loss_for_matched_load_at_frequency := CalculateFeedlineLossForMatchedLoadAtFrequency(feedline_length, feedline_loss_per_100ft_at_frequency)

	feedline_loss_for_matched_load_at_frequency_percentage := CalculateFeedlineLossForMatchedLoadAtFrequencyPercentage(feedline_loss_for_matched_load_at_frequency)

	gamma_squared := math.Pow(math.Abs(gamma), 2)

	feedline_loss_for_swr := CalculateFeedlineLossForSWR(feedline_loss_for_matched_load_at_frequency_percentage, gamma_squared)

	feedline_loss_for_swr_percentage := CalculateFeedlineLossForSWRPercentage(feedline_loss_for_swr)

	power_loss_at_swr := feedline_loss_for_swr_percentage * float64(transmitter_power)

	peak_envelope_power_at_antenna := float64(transmitter_power) - power_loss_at_swr

	uncontrolled_average_pep := peak_envelope_power_at_antenna * duty_cycle * uncontrolled_percentage_30_minutes

	mpe_s := 180 / (math.Pow(freq_values.freq, 2))

	gain_decimal := math.Pow(10, freq_values.gaindbi/10)

	return math.Sqrt((0.219 * uncontrolled_average_pep * gain_decimal) / mpe_s)
}

func CalculateReflectionCoefficient(freq_values FrequencyValues) float64 {
	return math.Abs(float64((freq_values.swr - 1) / (freq_values.swr + 1)))
}

func CalculateFeedlineLossForMatchedLoadAtFrequency(feedline_length int16, feedline_loss_per_100ft_at_frequency float64) float64 {
	return ((float64(feedline_length) / 100.0) * feedline_loss_per_100ft_at_frequency)
}

func CalculateFeedlineLossForMatchedLoadAtFrequencyPercentage(feedline_loss_for_matched_load float64) float64 {
	return math.Pow(10, (-(feedline_loss_for_matched_load) / 10.0))
}

func CalculateFeedlineLossPer100ftAtFrequency(freq_values FrequencyValues, cable_values CableValues) float64 {
	return cable_values.k1 * (math.Sqrt(freq_values.freq + cable_values.k2*freq_values.freq))
}

func CalculateFeedlineLossForSWR(feedline_loss_for_matched_load_percentage float64, gamma_squared float64) float64 {
	return -10 * math.Log10(feedline_loss_for_matched_load_percentage*
		((1-gamma_squared)/(1-math.Pow(feedline_loss_for_matched_load_percentage, 2)*gamma_squared)))
}

func CalculateFeedlineLossForSWRPercentage(feedline_loss_for_swr float64) float64 {
	return (100 - 100/(math.Pow(10, feedline_loss_for_swr/10))) / 100
}

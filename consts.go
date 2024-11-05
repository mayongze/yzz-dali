package homeassistant

type QueryColourValueDTR byte

func (qcvd QueryColourValueDTR) Byte() byte {
	return byte(qcvd)
}

const (
	XCoordinate                        QueryColourValueDTR = 0
	YCoordinate                        QueryColourValueDTR = 1
	ColourTemperatureTC                QueryColourValueDTR = 2
	PrimaryNDimLevel0                  QueryColourValueDTR = 3
	PrimaryNDimLevel1                  QueryColourValueDTR = 4
	PrimaryNDimLevel2                  QueryColourValueDTR = 5
	PrimaryNDimLevel3                  QueryColourValueDTR = 6
	PrimaryNDimLevel4                  QueryColourValueDTR = 7
	PrimaryNDimLevel5                  QueryColourValueDTR = 8
	RedDimLevel                        QueryColourValueDTR = 9
	GreenDimLevel                      QueryColourValueDTR = 10
	BlueDimLevel                       QueryColourValueDTR = 11
	WhiteDimLevel                      QueryColourValueDTR = 12
	AmberDimLevel                      QueryColourValueDTR = 13
	FreecolourDimLevel                 QueryColourValueDTR = 14
	RGBWAFControl                      QueryColourValueDTR = 15
	XCoordinatePrimaryN0               QueryColourValueDTR = 64
	YCoordinatePrimaryN0               QueryColourValueDTR = 65
	TYPrimaryN0                        QueryColourValueDTR = 66
	XCoordinatePrimaryN1               QueryColourValueDTR = 67
	YCoordinatePrimaryN1               QueryColourValueDTR = 68
	TYPrimaryN1                        QueryColourValueDTR = 69
	XCoordinatePrimaryN2               QueryColourValueDTR = 70
	YCoordinatePrimaryN2               QueryColourValueDTR = 71
	TYPrimaryN2                        QueryColourValueDTR = 72
	XCoordinatePrimaryN3               QueryColourValueDTR = 73
	YCoordinatePrimaryN3               QueryColourValueDTR = 74
	TYPrimaryN3                        QueryColourValueDTR = 75
	XCoordinatePrimaryN4               QueryColourValueDTR = 76
	YCoordinatePrimaryN4               QueryColourValueDTR = 77
	TYPrimaryN4                        QueryColourValueDTR = 78
	XCoordinatePrimaryN5               QueryColourValueDTR = 79
	YCoordinatePrimaryN5               QueryColourValueDTR = 80
	TYPrimaryN5                        QueryColourValueDTR = 81
	NumberOfPrimaries                  QueryColourValueDTR = 82
	ColourTemperatureTcCoolest         QueryColourValueDTR = 128
	ColourTemperatureTcPhysicalCoolest QueryColourValueDTR = 129
	ColourTemperatureTcWarmest         QueryColourValueDTR = 130
	ColourTemperatureTcPhysicalWarmest QueryColourValueDTR = 131
	TemporaryXCoordinate               QueryColourValueDTR = 192
	TemporaryYCoordinate               QueryColourValueDTR = 193
	TemporaryColourTemperature         QueryColourValueDTR = 194
	TemporaryPrimaryNDimLevel0         QueryColourValueDTR = 195
	TemporaryPrimaryNDimLevel1         QueryColourValueDTR = 196
	TemporaryPrimaryNDimLevel2         QueryColourValueDTR = 197
	TemporaryPrimaryNDimLevel3         QueryColourValueDTR = 198
	TemporaryPrimaryNDimLevel4         QueryColourValueDTR = 199
	TemporaryPrimaryNDimLevel5         QueryColourValueDTR = 200
	TemporaryRedDimLevel               QueryColourValueDTR = 201
	TemporaryGreenDimLevel             QueryColourValueDTR = 202
	TemporaryBlueDimLevel              QueryColourValueDTR = 203
	TemporaryWhiteDimLevel             QueryColourValueDTR = 204
	TemporaryAmberDimLevel             QueryColourValueDTR = 205
	TemporaryFreecolourDimLevel        QueryColourValueDTR = 206
	TemporaryRgbwafControl             QueryColourValueDTR = 207
	TemporaryColourType                QueryColourValueDTR = 208
	ReportXCoordinate                  QueryColourValueDTR = 224
	ReportYCoordinate                  QueryColourValueDTR = 225
	ReportColourTemperatureTc          QueryColourValueDTR = 226
	ReportPrimaryNDimLevel0            QueryColourValueDTR = 227
	ReportPrimaryNDimLevel1            QueryColourValueDTR = 228
	ReportPrimaryNDimLevel2            QueryColourValueDTR = 229
	ReportPrimaryNDimLevel3            QueryColourValueDTR = 230
	ReportPrimaryNDimLevel4            QueryColourValueDTR = 231
	ReportPrimaryNDimLevel5            QueryColourValueDTR = 232
	ReportRedDimLevel                  QueryColourValueDTR = 233
	ReportGreenDimLevel                QueryColourValueDTR = 234
	ReportBlueDimLevel                 QueryColourValueDTR = 235
	ReportWhiteDimLevel                QueryColourValueDTR = 236
	ReportAmberDimLevel                QueryColourValueDTR = 237
	ReportFreecolourDimLevel           QueryColourValueDTR = 238
	ReportRgbwafControl                QueryColourValueDTR = 239
	ReportColourType                   QueryColourValueDTR = 240
)

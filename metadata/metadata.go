package metadata

type MetaData struct {
	mSource      string
	mDestination string
	mAssetName   string
}

func CreateNewMetaData(source string, destination string, assetName string) *MetaData {
	return &MetaData{
		mSource:      source,
		mDestination: destination,
		mAssetName:   assetName,
	}
}

func (md *MetaData) GetSourceName() string {
	return md.mSource
}

func (md *MetaData) GetDestinationName() string {
	return md.mDestination
}

func (md *MetaData) GetAssetName() string {
	return md.mAssetName
}

func (md *MetaData) GetMetaDataString() string {
	return md.mSource + " " + md.mDestination + " " + md.mAssetName
}

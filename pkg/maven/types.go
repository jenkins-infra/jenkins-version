package maven

// Metadata struct to store metadata.
type Metadata struct {
	GroupID    string     `xml:"groupId"`
	ArtifactID string     `xml:"artifactId"`
	Versioning Versioning `xml:"versioning"`
}

// Versioning struct to store versioning.
type Versioning struct {
	Latest   string   `xml:"latest"`
	Release  string   `xml:"release"`
	Versions Versions `xml:"versions"`
}

// Versions struct to store versions.
type Versions struct {
	Versions []string `xml:"version"`
}

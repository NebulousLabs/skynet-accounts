package test

import (
	"context"
	"testing"

	"github.com/NebulousLabs/skynet-accounts/database"
	"github.com/NebulousLabs/skynet-accounts/skynet"

	"gitlab.com/NebulousLabs/fastrand"
)

// TestUserStats ensures we report accurate statistics for users.
func TestUserStats(t *testing.T) {
	ctx := context.Background()
	db, err := database.New(ctx, DBTestCredentials(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add a test user.
	sub := string(fastrand.Bytes(userSubLen))
	u, err := db.UserCreate(nil, sub, database.TierPremium5)
	if err != nil {
		t.Fatal(err)
	}
	defer func(user *database.User) {
		_ = db.UserDelete(nil, user)
	}(u)

	testUploadSizeSmall := int64(1 + fastrand.Intn(4*skynet.MiB-1))
	testUploadSizeBig := int64(4*skynet.MiB + 1 + fastrand.Intn(4*skynet.MiB))
	expectedUploadBandwidth := int64(0)
	expectedDownloadBandwidth := int64(0)

	// Create a small upload.
	skylinkSmall, err := createTestUpload(ctx, db, u, testUploadSizeSmall)
	if err != nil {
		t.Fatal(err)
	}
	expectedUploadBandwidth = skynet.BandwidthUploadCost(testUploadSizeSmall)
	// Check the stats.
	stats, err := db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user stats.", err)
	}
	if stats.NumUploads != 1 {
		t.Fatalf("Expected a total of %d uploads, got %d.", 1, stats.NumUploads)
	}
	if stats.BandwidthUploads != expectedUploadBandwidth {
		t.Fatalf("Expected upload bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedUploadBandwidth, expectedUploadBandwidth/skynet.MiB,
			stats.BandwidthUploads, stats.BandwidthUploads/skynet.MiB)
	}

	// Create a big upload.
	skylinkBig, err := createTestUpload(ctx, db, u, testUploadSizeBig)
	if err != nil {
		t.Fatal(err)
	}
	expectedUploadBandwidth += skynet.BandwidthUploadCost(testUploadSizeBig)
	// Check the stats.
	stats, err = db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user stats.", err)
	}
	if stats.NumUploads != 2 {
		t.Fatalf("Expected a total of %d uploads, got %d.", 2, stats.NumUploads)
	}
	if stats.BandwidthUploads != expectedUploadBandwidth {
		t.Fatalf("Expected upload bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedUploadBandwidth, expectedUploadBandwidth/skynet.MiB,
			stats.BandwidthUploads, stats.BandwidthUploads/skynet.MiB)
	}

	// Register a small download.
	smallDownload := int64(1 + fastrand.Intn(4*skynet.MiB))
	err = db.DownloadCreate(ctx, *u, *skylinkSmall, smallDownload)
	if err != nil {
		t.Fatal("Failed to download.", err)
	}
	expectedDownloadBandwidth += skynet.BandwidthDownloadCost(smallDownload)
	// Check the stats.
	stats, err = db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user stats.", err)
	}
	if stats.NumDownloads != 1 {
		t.Fatalf("Expected a total of %d downloads, got %d.", 1, stats.NumDownloads)
	}
	if stats.BandwidthDownloads != expectedDownloadBandwidth {
		t.Fatalf("Expected download bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedDownloadBandwidth, expectedDownloadBandwidth/skynet.MiB,
			stats.BandwidthDownloads, stats.BandwidthDownloads/skynet.MiB)
	}
	// Register a big download.
	bigDownload := int64(100*skynet.MiB + fastrand.Intn(4*skynet.MiB))
	err = db.DownloadCreate(ctx, *u, *skylinkBig, bigDownload)
	if err != nil {
		t.Fatal("Failed to download.", err)
	}
	expectedDownloadBandwidth += skynet.BandwidthDownloadCost(bigDownload)
	// Check bandwidth.
	stats, err = db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user stats.", err)
	}
	if stats.NumDownloads != 2 {
		t.Fatalf("Expected a total of %d downloads, got %d.", 2, stats.NumDownloads)
	}
	if stats.BandwidthDownloads != expectedDownloadBandwidth {
		t.Fatalf("Expected download bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedDownloadBandwidth, expectedDownloadBandwidth/skynet.MiB,
			stats.BandwidthDownloads, stats.BandwidthDownloads/skynet.MiB)
	}

	// Register a registry read.
	_, err = db.RegistryReadCreate(ctx, *u)
	if err != nil {
		t.Fatal("Failed to register a registry read.", err)
	}
	expectedRegReadBandwidth := int64(skynet.PriceBandwidthRegistryRead)
	// Check bandwidth.
	stats, err = db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user details.", err)
	}
	if stats.NumRegReads != 1 {
		t.Fatalf("Expected a total of %d registry reads, got %d.", 1, stats.NumRegReads)
	}
	if stats.BandwidthRegReads != expectedRegReadBandwidth {
		t.Fatalf("Expected registry read bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedRegReadBandwidth, expectedRegReadBandwidth/skynet.MiB,
			stats.BandwidthRegReads, stats.BandwidthRegReads/skynet.MiB)
	}
	// Register a registry read.
	_, err = db.RegistryReadCreate(ctx, *u)
	if err != nil {
		t.Fatal("Failed to register a registry read.", err)
	}
	expectedRegReadBandwidth += int64(skynet.PriceBandwidthRegistryRead)
	// Check bandwidth.
	stats, err = db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user details.", err)
	}
	if stats.NumRegReads != 2 {
		t.Fatalf("Expected a total of %d registry reads, got %d.", 2, stats.NumRegReads)
	}
	if stats.BandwidthRegReads != expectedRegReadBandwidth {
		t.Fatalf("Expected registry read bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedRegReadBandwidth, expectedRegReadBandwidth/skynet.MiB,
			stats.BandwidthRegReads, stats.BandwidthRegReads/skynet.MiB)
	}

	// Register a registry write.
	_, err = db.RegistryWriteCreate(ctx, *u)
	if err != nil {
		t.Fatal("Failed to register a registry write.", err)
	}
	expectedRegWriteBandwidth := int64(skynet.PriceBandwidthRegistryWrite)
	// Check bandwidth.
	stats, err = db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user details.", err)
	}
	if stats.NumRegWrites != 1 {
		t.Fatalf("Expected a total of %d registry writes, got %d.", 1, stats.NumRegWrites)
	}
	if stats.BandwidthRegWrites != expectedRegWriteBandwidth {
		t.Fatalf("Expected registry write bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedRegWriteBandwidth, expectedRegWriteBandwidth/skynet.MiB,
			stats.BandwidthRegWrites, stats.BandwidthRegWrites/skynet.MiB)
	}
	// Register a registry write.
	_, err = db.RegistryWriteCreate(ctx, *u)
	if err != nil {
		t.Fatal("Failed to register a registry write.", err)
	}
	expectedRegWriteBandwidth += int64(skynet.PriceBandwidthRegistryWrite)
	// Check bandwidth.
	stats, err = db.UserStats(ctx, *u)
	if err != nil {
		t.Fatal("Failed to fetch user details.", err)
	}
	if stats.NumRegWrites != 2 {
		t.Fatalf("Expected a total of %d registry writes, got %d.", 2, stats.NumRegWrites)
	}
	if stats.BandwidthRegWrites != expectedRegWriteBandwidth {
		t.Fatalf("Expected registry write bandwidth of %d (%d MiB), got %d (%d MiB).",
			expectedRegWriteBandwidth, expectedRegWriteBandwidth/skynet.MiB,
			stats.BandwidthRegWrites, stats.BandwidthRegWrites/skynet.MiB)
	}
}

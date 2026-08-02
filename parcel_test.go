package main

import (
	"database/sql"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "tracker.db")

	db, err := sql.Open("sqlite", dbPath)
	require.NoError(t, err)

	err = createSchema(db)
	require.NoError(t, err)

	return db
}

func TestAddGetDelete(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	parcel.Number = id

	got, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, parcel, got)

	err = store.Delete(id)
	require.NoError(t, err)

	_, err = store.Get(id)
	require.Error(t, err)
}

func TestSetAddress(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	newAddress := "new test address"

	err = store.SetAddress(id, newAddress)
	require.NoError(t, err)

	got, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, newAddress, got.Address)

	err = store.Delete(id)
	require.NoError(t, err)
}

func TestSetStatus(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	id, err := store.Add(parcel)
	require.NoError(t, err)
	require.NotZero(t, id)

	err = store.SetStatus(id, ParcelStatusSent)
	require.NoError(t, err)

	got, err := store.Get(id)
	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, got.Status)
}

func TestGetByClient(t *testing.T) {
	db := openTestDB(t)
	defer db.Close()

	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := 42

	for i := 0; i < len(parcels); i++ {
		parcels[i].Client = client

		id, err := store.Add(parcels[i])
		require.NoError(t, err)
		require.NotZero(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		expected, ok := parcelMap[parcel.Number]
		require.True(t, ok)
		require.Equal(t, expected, parcel)
	}
}

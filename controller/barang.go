package controller

import (
	"crud-go/config"
	"crud-go/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetAllBarang - Get all barang with pagination
func GetAllBarang(c *gin.Context) {
	var barangs []model.Barang

	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "5")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		limitInt = 5
	}

	offset := (pageInt - 1) * limitInt

	var total int
	err = config.DB.QueryRow("SELECT COUNT(*) FROM barang").Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}

	rows, err := config.DB.Query("SELECT id, nama_barang, deskripsi, kategori, harga, stok, gambar, createdAt, updatedAt FROM barang LIMIT ? OFFSET ?", limitInt, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"result": err.Error(),
			"status": 0,
		})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var b model.Barang
		err := rows.Scan(
			&b.Id,
			&b.Nama_barang,
			&b.Deskripsi,
			&b.Kategori,
			&b.Harga,
			&b.Stok,
			&b.Gambar,
			&b.Created_at,
			&b.Updated_at,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"result": err.Error(),
				"status": 0,
			})
			return
		}
		barangs = append(barangs, b)
	}

	c.JSON(http.StatusOK, gin.H{
		"result": barangs,
		"status": 1,
		"page":   pageInt,
		"limit":  limitInt,
		"total":  total,
	})
}

// GetBarangByID - Get barang by ID
func GetBarangByID(c *gin.Context) {
	id := c.Param("id")
	var barang model.Barang

	query := "SELECT id, nama_barang, deskripsi, kategori, harga, stok, gambar, createdAt, updatedAt FROM barang WHERE id = ?"
	err := config.DB.QueryRow(query, id).Scan(
		&barang.Id,
		&barang.Nama_barang,
		&barang.Deskripsi,
		&barang.Kategori,
		&barang.Harga,
		&barang.Stok,
		&barang.Gambar,
		&barang.Created_at,
		&barang.Updated_at,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Barang tidak ditemukan",
			"status": 0,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": barang,
		"status": 1,
	})
}

// CreateBarang - Create new barang (Admin & Seller only)
func CreateBarang(c *gin.Context) {
	// Check role - only admin and seller can create barang
	userRole := GetUserRoleFromContext(c)
	if userRole != "admin" && userRole != "seller" {
		c.JSON(http.StatusForbidden, gin.H{
			"err":    "Hanya admin dan seller yang dapat menambahkan barang",
			"status": 0,
		})
		return
	}

	var newBarang model.Barang

	err := c.BindJSON(&newBarang)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	newBarang.Created_at = time.Now()
	newBarang.Updated_at = time.Now()

	query := "INSERT INTO barang (nama_barang, deskripsi, kategori, harga, stok, gambar, createdAt, updatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := config.DB.Exec(query, newBarang.Nama_barang, newBarang.Deskripsi, newBarang.Kategori, newBarang.Harga, newBarang.Stok, newBarang.Gambar, newBarang.Created_at, newBarang.Updated_at)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "Barang berhasil ditambahkan",
			"status": 1,
			"data":   newBarang,
		})
	}
}

// UpdateBarang - Update barang by ID (Admin only)
func UpdateBarang(c *gin.Context) {
	// Check role - only admin can update barang
	userRole := GetUserRoleFromContext(c)
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"err":    "Hanya admin yang dapat mengupdate barang",
			"status": 0,
		})
		return
	}

	id := c.Param("id")
	var updatedBarang model.Barang

	err := c.BindJSON(&updatedBarang)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	updatedBarang.Updated_at = time.Now()

	query := "UPDATE barang SET nama_barang = ?, deskripsi = ?, kategori = ?, harga = ?, stok = ?, gambar = ?, updatedAt = ? WHERE id = ?"
	result, err := config.DB.Exec(query, updatedBarang.Nama_barang, updatedBarang.Deskripsi, updatedBarang.Kategori, updatedBarang.Harga, updatedBarang.Stok, updatedBarang.Gambar, updatedBarang.Updated_at, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "Barang berhasil diperbarui",
			"status": 1,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Barang tidak ditemukan",
			"status": 0,
		})
	}
}

// DeleteBarang - Delete barang by ID (Admin & Seller)
func DeleteBarang(c *gin.Context) {
	// Check role - only admin and seller can delete barang
	userRole := GetUserRoleFromContext(c)
	if userRole != "admin" && userRole != "seller" {
		c.JSON(http.StatusForbidden, gin.H{
			"err":    "Hanya admin dan seller yang dapat menghapus barang",
			"status": 0,
		})
		return
	}

	id := c.Param("id")

	query := "DELETE FROM barang WHERE id = ?"
	result, err := config.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err":    err.Error(),
			"status": 0,
		})
		return
	}

	if affected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"result": "Barang berhasil dihapus",
			"status": 1,
		})
	} else {
		c.JSON(http.StatusNotFound, gin.H{
			"result": "Barang tidak ditemukan",
			"status": 0,
		})
	}
}

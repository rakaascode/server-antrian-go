-- ============================================================
-- SCRIPT FIX: Reset password admin cabang ke bcrypt hash
-- Jalankan di psql setelah server di-restart
-- ============================================================

-- Lihat semua admin dan status passwordnya
SELECT id, name, username, role, cabang_id,
       CASE WHEN password LIKE '$2a$%' THEN 'bcrypt ✅' ELSE 'PLAINTEXT ❌' END AS password_status
FROM users
WHERE role = 'admin'
ORDER BY id;

-- ============================================================
-- OPSI A: Reset semua admin cabang (non-seeder) ke password baru
-- Password baru: Admin@Reset123
-- bcrypt hash dari "Admin@Reset123":
-- =====================================d=======================
-- Hasilkan hash dengan Go (jalankan di terminal terpisah):
--   go run -e 'import "golang.org/x/crypto/bcrypt"; h,_:=bcrypt.GenerateFromPassword([]byte("Admin@Reset123"),bcrypt.DefaultCost); fmt.Println(string(h))'

-- ATAU gunakan endpoint PUT /api/users/:id setelah server restart
-- (password otomatis di-hash oleh service layer yang sudah difix)

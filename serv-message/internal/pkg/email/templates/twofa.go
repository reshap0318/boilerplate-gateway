package templates

import "fmt"

// TwoFACodeContent renders the content section of the 2FA login-code
// email. Wrap it with Layout before sending.
func TwoFACodeContent(code string) string {
	return fmt.Sprintf(`
<h2>Kode Verifikasi Login</h2>
<p>Halo,</p>
<p>Gunakan kode berikut untuk menyelesaikan proses login Anda:</p>

<p class="token-box" style="text-align: center; letter-spacing: 4px; font-weight: 600; font-size: 24px;">%s</p>

<div class="warning">
    <strong>Penting:</strong> Kode ini hanya berlaku selama 5 menit. Jika Anda tidak mencoba login, abaikan email ini atau segera hubungi layanan dukungan kami.
</div>
`, code)
}

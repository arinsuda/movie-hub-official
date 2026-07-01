<template>
  <div class="donate-page">
    <section class="hero">
      <span class="eyebrow">สนับสนุนผู้พัฒนา</span>
      <h1>ขอบคุณที่อยากซัพพอร์ต 🙏</h1>
      <p>
        <!-- TODO: แก้ข้อความตามที่ต้องการ -->
        REMOV เป็นโปรเจกต์ที่พัฒนาโดย [ชื่อของคุณ] ทุกการสนับสนุน
        ไม่ว่าเล็กหรือใหญ่ ช่วยให้ผมมีแรงพัฒนาฟีเจอร์ใหม่ๆ ต่อไปครับ
      </p>
    </section>

    <div class="donate-grid">
      <!-- PromptPay QR -->
      <div class="donate-card">
        <div class="card-head">
          <QrCode :size="18" />
          <h3>PromptPay</h3>
        </div>

        <div class="qr-box">
          <!-- TODO: เปลี่ยนเป็นไฟล์ QR PromptPay จริงของคุณ -->
          <img
            src="/public/mePromptPay.png"
            alt="PromptPay QR Code"
            @error="onQrError"
          />
          <p v-if="qrMissing" class="qr-fallback">
            วางไฟล์ QR ไว้ที่<br />
            <code>client/public/images/promptpay-qr.png</code>
          </p>
        </div>

        <p class="qr-name">
          <!-- TODO: ใส่ชื่อบัญชี PromptPay -->
          ชื่อบัญชี: นาย อรินทร์ สุดากิจจาทร
        </p>
        <p class="qr-hint">สแกนด้วยแอปธนาคารเพื่อโอนผ่าน PromptPay</p>
      </div>

      <!-- Bank Transfer -->
      <div class="donate-card">
        <div class="card-head">
          <Landmark :size="18" />
          <h3>โอนผ่านธนาคาร</h3>
        </div>

        <div class="bank-row">
          <!-- TODO: เปลี่ยนโลโก้/ชื่อธนาคาร -->
          <div class="bank-logo">
            <img
              :src="bank.logo"
              :alt="`${bank.name} logo`"
              @error="onBankLogoError"
            />
            <span v-if="bankLogoMissing">{{ bank.shortName }}</span>
          </div>
          <div class="bank-detail">
            <span class="bank-name">{{ bank.name }}</span>
            <span class="bank-account-name">{{ bank.accountName }}</span>
          </div>
        </div>

        <div class="account-box">
          <span class="account-number">{{ formattedAccount }}</span>
          <button class="copy-btn" @click="copyAccount">
            <component :is="copied ? Check : Copy" :size="14" />
            {{ copied ? "คัดลอกแล้ว" : "คัดลอก" }}
          </button>
        </div>
      </div>
    </div>

    <!-- Optional: other platforms -->
    <!-- <section class="alt-support">
      <p class="alt-label">หรือสนับสนุนผ่านช่องทางอื่น</p>
      <div class="alt-links">
        <a href="https://ko-fi.com/USERNAME" target="_blank" class="alt-btn">
          <Coffee :size="15" /> Ko-fi
        </a>
        <a
          href="https://www.buymeacoffee.com/USERNAME"
          target="_blank"
          class="alt-btn"
        >
          <Coffee :size="15" /> Buy Me a Coffee
        </a>
      </div>
    </section> -->

    <p class="thanks-note">
      💛 ขอบคุณทุกกำลังใจและการสนับสนุน ไม่ว่าคุณจะโดเนทหรือแค่เข้ามาใช้งาน
      REMOV ก็มีความหมายกับผมมากครับ
    </p>
  </div>
</template>

<script setup lang="ts">
  import { computed, ref } from "vue"
  import { QrCode, Landmark, Copy, Check, Coffee } from "lucide-vue-next"

  const bank = {
    shortName: "KBANK",
    name: "ธนาคารกสิกรไทย",
    accountName: "นาย อรินทร์ สุดากิจจาทร",
    accountNumber: "0453315436",
    logo: "/public/kbank-logo.png",
  }

  const bankLogoMissing = ref(false)

  function onBankLogoError() {
    bankLogoMissing.value = true
  }

  const formattedAccount = computed(() => bank.accountNumber)

  const copied = ref(false)
  async function copyAccount() {
    try {
      await navigator.clipboard.writeText(bank.accountNumber.replace(/-/g, ""))
      copied.value = true
      setTimeout(() => (copied.value = false), 2000)
    } catch (err) {
      console.error("Copy failed:", err)
    }
  }

  const qrMissing = ref(false)
  function onQrError() {
    qrMissing.value = true
  }
</script>

<style scoped>
  .donate-page {
    --c-bg: #101010;
    --c-card: #161616;
    --c-border: rgba(255, 255, 255, 0.06);
    --c-border-h: rgba(255, 255, 255, 0.12);
    --c-red: #e1251b;
    --c-text: #f0f0f0;
    --c-sub: #8a8a8e;
    --c-muted: #3a3a3c;
    --ease: cubic-bezier(0.16, 1, 0.3, 1);
    color: var(--c-text);
    max-width: 720px;
    margin: 0 auto;
    padding: 48px 24px 80px;
  }

  .hero {
    text-align: center;
    padding-bottom: 40px;
    border-bottom: 1px solid var(--c-border);
    margin-bottom: 40px;
  }
  .eyebrow {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--c-red);
  }
  .hero h1 {
    font-size: 1.8rem;
    font-weight: 800;
    margin: 10px 0 12px;
  }
  .hero p {
    color: var(--c-sub);
    font-size: 0.88rem;
    line-height: 1.7;
    max-width: 480px;
    margin: 0 auto;
  }

  .donate-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
    gap: 16px;
    margin-bottom: 40px;
  }

  .donate-card {
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 16px;
    padding: 22px;
  }
  .card-head {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 16px;
    color: var(--c-red);
  }
  .card-head h3 {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--c-text);
    margin: 0;
  }

  /* QR */
  .qr-box {
    width: 100%;
    aspect-ratio: 1;
    background: #fff;
    border-radius: 12px;
    overflow: hidden;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 14px;
    position: relative;
  }
  .qr-box img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
  .qr-fallback {
    position: absolute;
    color: #333;
    font-size: 0.72rem;
    text-align: center;
    padding: 12px;
    line-height: 1.6;
  }
  .qr-fallback code {
    font-size: 0.65rem;
    background: rgba(0, 0, 0, 0.06);
    padding: 2px 5px;
    border-radius: 4px;
  }
  .qr-name {
    font-size: 0.82rem;
    font-weight: 600;
    margin: 0 0 4px;
  }
  .qr-hint {
    font-size: 0.72rem;
    color: var(--c-sub);
    margin: 0;
  }

  /* Bank */
  .bank-row {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 18px;
  }
  .bank-logo {
    width: 44px;
    height: 44px;
    flex-shrink: 0;
    border-radius: 10px;
    background: #0d0d0d;
    border: 1px solid var(--c-border);
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: hidden;
    padding: 6px;
  }

  .bank-logo img {
    width: 100%;
    height: 100%;
    object-fit: contain;
    display: block;
  }

  .bank-logo span {
    font-size: 0.6rem;
    font-weight: 800;
    color: var(--c-red);
  }
  .bank-detail {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }
  .bank-name {
    font-size: 0.84rem;
    font-weight: 600;
  }
  .bank-account-name {
    font-size: 0.74rem;
    color: var(--c-sub);
  }

  .account-box {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: #0d0d0d;
    border: 1px solid var(--c-border);
    border-radius: 10px;
    padding: 12px 14px;
  }
  .account-number {
    font-size: 0.95rem;
    font-weight: 700;
    font-variant-numeric: tabular-nums;
    letter-spacing: 0.03em;
  }
  .copy-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid var(--c-border-h);
    color: var(--c-text);
    font-size: 0.72rem;
    font-weight: 600;
    padding: 6px 10px;
    border-radius: 7px;
    cursor: pointer;
    transition: all 0.15s;
  }
  .copy-btn:hover {
    background: rgba(255, 255, 255, 0.1);
  }

  /* Alt support */
  .alt-support {
    text-align: center;
    margin-bottom: 40px;
  }
  .alt-label {
    font-size: 0.76rem;
    color: var(--c-sub);
    margin-bottom: 12px;
  }
  .alt-links {
    display: flex;
    justify-content: center;
    gap: 10px;
    flex-wrap: wrap;
  }
  .alt-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: var(--c-card);
    border: 1px solid var(--c-border);
    color: var(--c-text);
    font-size: 0.78rem;
    font-weight: 600;
    padding: 8px 16px;
    border-radius: 8px;
    text-decoration: none;
    transition: border-color 0.15s;
  }
  .alt-btn:hover {
    border-color: var(--c-border-h);
  }

  .thanks-note {
    text-align: center;
    font-size: 0.8rem;
    color: var(--c-sub);
    line-height: 1.7;
    max-width: 480px;
    margin: 0 auto;
  }
</style>

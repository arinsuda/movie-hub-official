<template>
  <div class="about-page">
    <!-- Hero -->
    <section class="hero">
      <span class="eyebrow">เกี่ยวกับเรา</span>
      <h1 class="hero-title">
        REMOV
        <span class="hero-accent"></span>
      </h1>
      <p class="hero-sub">
        พื้นที่สำหรับคนรักหนังและซีรีส์ ค้นหา รีวิว
        และติดตามสิ่งที่คุณดูในที่เดียว
      </p>
    </section>

    <!-- Mission -->
    <section class="block">
      <div class="block-head">
        <span class="tag">01</span>
        <h2>ทำไมถึงสร้าง REMOV</h2>
      </div>
      <p class="block-text">
        REMOV เริ่มต้นจากความชอบในการดูภาพยนตร์ของผมเอง
        ทุกครั้งที่ได้ดูหนังที่ชอบ ผมมักอยากแชร์ความรู้สึก ความคิดเห็น
        หรือเล่าเรื่องราวให้คนอื่นที่สนใจเหมือนกันฟัง แต่ปัญหาคือ
        เรามักไม่รู้ว่าเพื่อนหรือคนรอบตัวชอบหนังแนวไหน
        ดูเรื่องเดียวกับเราหรือไม่ หรือมีรสนิยมใกล้เคียงกันแค่ไหน
        จนกว่าจะได้พูดคุยกันหลายครั้ง ผมจึงอยากสร้าง REMOV
        ให้เป็นแพลตฟอร์มสำหรับคนรักหนัง ที่ช่วยให้ผู้ใช้สามารถบันทึกสิ่งที่ดู
        เขียนรีวิว แบ่งปันความคิดเห็น และค้นพบภาพยนตร์หรือซีรีส์เรื่องใหม่ ๆ
        ได้ง่ายขึ้น พร้อมเชื่อมต่อกับผู้คนที่มีความสนใจใกล้เคียงกัน
      </p>
    </section>

    <!-- Features -->
    <section class="block">
      <div class="block-head">
        <span class="tag">02</span>
        <h2>สิ่งที่ทำได้บน REMOV</h2>
      </div>
      <div class="feature-grid">
        <div v-for="f in features" :key="f.title" class="feature-card">
          <component :is="f.icon" :size="20" :stroke-width="1.6" />
          <h3>{{ f.title }}</h3>
          <p>{{ f.desc }}</p>
        </div>
      </div>
    </section>

    <!-- Tech stack -->
    <section class="block">
      <div class="block-head">
        <span class="tag">03</span>
        <h2>เบื้องหลังการพัฒนา</h2>
      </div>
      <p class="block-text">
        REMOV พัฒนาโดยใช้ Go (Fiber, GORM, PostgreSQL) ฝั่ง Backend และ Vue 3
        (Composition API, TypeScript, Pinia, TanStack Query) ฝั่ง Frontend
      </p>
      <div class="stack-row">
        <span v-for="s in stack" :key="s" class="stack-chip">{{ s }}</span>
      </div>
    </section>

    <!-- Creator -->
    <section class="creator-card">
      <div class="creator-avatar">
        <!-- TODO: เปลี่ยนเป็นรูปโปรไฟล์จริง -->
        <img src="/public/arin.png" alt="ผู้พัฒนา" @error="onAvatarError" />
      </div>
      <div class="creator-info">
        <span class="eyebrow">ผู้พัฒนา</span>
        <!-- TODO: ใส่ชื่อจริง/ชื่อเล่น -->
        <h3>อรินทร์ สุดากิจจาทร</h3>
        <p>
          <!-- TODO: bio สั้นๆ -->
          นักศึกษาด้าน IT ที่หลงรักการสร้างโปรดักต์ Full-stack และภาพยนตร์
        </p>
        <div class="creator-links">
          <!-- TODO: แก้ลิงก์ทั้งหมด -->
          <a href="https://github.com/arinsuda" target="_blank">
            <Github :size="15" /> GitHub
          </a>
          <a href="https://linkedin.com/in/arinsuda" target="_blank">
            <Linkedin :size="15" /> LinkedIn
          </a>
          <a href="mailto:sixarin.thorn@gmail.com">
            <Mail :size="15" /> Email
          </a>
        </div>
      </div>
    </section>

    <!-- CTA to donate -->
    <section class="support-cta">
      <div>
        <h3>ชอบ REMOV ไหม?</h3>
        <p>ถ้าอยากสนับสนุนให้โปรเจกต์นี้พัฒนาต่อไป กดปุ่มด้านขวาได้เลย</p>
      </div>
      <RouterLink :to="{ name: 'donate' }" class="support-btn">
        <Heart :size="15" /> สนับสนุนผู้พัฒนา
      </RouterLink>
    </section>
  </div>
</template>

<script setup lang="ts">
  import {
    Star,
    Bookmark,
    Trophy,
    Users,
    Github,
    Linkedin,
    Mail,
    Heart,
  } from "lucide-vue-next"

  const features = [
    {
      icon: Star,
      title: "รีวิว & เรตติ้ง",
      desc: "เขียนรีวิวและให้คะแนนหนัง/ซีรีส์ที่คุณดู",
    },
    {
      icon: Bookmark,
      title: "Watchlist",
      desc: "บันทึกรายการที่อยากดูและติดตามสิ่งที่ดูแล้ว",
    },
    {
      icon: Trophy,
      title: "Achievement",
      desc: "ปลดล็อกความสำเร็จตามพฤติกรรมการดูของคุณ",
    },
    {
      icon: Users,
      title: "Community",
      desc: "ติดตามผู้ใช้คนอื่นและดูรีวิวจากฟีดของคุณ",
    },
  ]

  const stack = [
    "Go",
    "Fiber v3",
    "GORM",
    "PostgreSQL",
    "Vue 3",
    "TypeScript",
    "Pinia",
    "TanStack Query",
  ]

  function onAvatarError(e: Event) {
    ;(e.target as HTMLImageElement).style.display = "none"
  }
</script>

<style scoped>
  .about-page {
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
    max-width: 900px;
    margin: 0 auto;
    padding: 48px 24px 80px;
  }

  .eyebrow {
    font-size: 0.65rem;
    font-weight: 700;
    letter-spacing: 0.12em;
    text-transform: uppercase;
    color: var(--c-red);
  }

  /* Hero */
  .hero {
    padding: 40px 0 56px;
    border-bottom: 1px solid var(--c-border);
    margin-bottom: 48px;
  }
  .hero-title {
    font-size: 2.6rem;
    font-weight: 800;
    margin: 10px 0 14px;
    letter-spacing: -0.02em;
  }
  .hero-accent {
    color: var(--c-red);
  }
  .hero-sub {
    color: var(--c-sub);
    font-size: 0.95rem;
    max-width: 480px;
    line-height: 1.6;
  }

  /* Blocks */
  .block {
    margin-bottom: 48px;
  }
  .block-head {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 14px;
  }
  .block-head h2 {
    font-size: 1.15rem;
    font-weight: 700;
    margin: 0;
  }
  .tag {
    font-size: 0.7rem;
    font-weight: 700;
    color: var(--c-muted);
    border: 1px solid var(--c-border);
    padding: 3px 8px;
    border-radius: 5px;
  }
  .block-text {
    color: var(--c-sub);
    font-size: 0.88rem;
    line-height: 1.7;
    max-width: 620px;
  }

  .feature-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 14px;
    margin-top: 16px;
  }
  .feature-card {
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 12px;
    padding: 18px;
    transition: border-color 0.2s var(--ease);
  }
  .feature-card:hover {
    border-color: var(--c-border-h);
  }
  .feature-card svg {
    color: var(--c-red);
    margin-bottom: 10px;
  }
  .feature-card h3 {
    font-size: 0.88rem;
    font-weight: 600;
    margin: 0 0 4px;
  }
  .feature-card p {
    font-size: 0.76rem;
    color: var(--c-sub);
    line-height: 1.5;
    margin: 0;
  }

  .stack-row {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 14px;
  }
  .stack-chip {
    font-size: 0.72rem;
    font-weight: 600;
    color: var(--c-sub);
    background: var(--c-card);
    border: 1px solid var(--c-border);
    padding: 5px 12px;
    border-radius: 20px;
  }

  /* Creator */
  .creator-card {
    display: flex;
    gap: 20px;
    align-items: flex-start;
    background: var(--c-card);
    border: 1px solid var(--c-border);
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 32px;
  }
  .creator-avatar {
    width: 72px;
    height: 72px;
    flex-shrink: 0;
    border-radius: 50%;
    overflow: hidden;
    background: #0d0d0d;
    border: 1px solid var(--c-border);
  }
  .creator-avatar img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
  .creator-info h3 {
    font-size: 1.1rem;
    font-weight: 700;
    margin: 4px 0 6px;
  }
  .creator-info p {
    font-size: 0.8rem;
    color: var(--c-sub);
    line-height: 1.6;
    margin: 0 0 12px;
  }
  .creator-links {
    display: flex;
    gap: 14px;
    flex-wrap: wrap;
  }
  .creator-links a {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    font-size: 0.76rem;
    font-weight: 600;
    color: var(--c-sub);
    text-decoration: none;
    transition: color 0.15s;
  }
  .creator-links a:hover {
    color: #fff;
  }

  /* CTA */
  .support-cta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;
    background: linear-gradient(135deg, rgba(225, 37, 27, 0.12), transparent);
    border: 1px solid var(--c-border);
    border-radius: 16px;
    padding: 24px;
  }
  .support-cta h3 {
    font-size: 1rem;
    margin: 0 0 4px;
  }
  .support-cta p {
    font-size: 0.8rem;
    color: var(--c-sub);
    margin: 0;
  }
  .support-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    background: var(--c-red);
    color: #fff;
    font-size: 0.82rem;
    font-weight: 700;
    padding: 10px 18px;
    border-radius: 8px;
    text-decoration: none;
    white-space: nowrap;
    transition: filter 0.15s;
  }
  .support-btn:hover {
    filter: brightness(1.1);
  }

  @media (max-width: 480px) {
    .creator-card {
      flex-direction: column;
      align-items: center;
      text-align: center;
    }
    .creator-links {
      justify-content: center;
    }
  }
</style>

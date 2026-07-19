import fs from "fs";
import path from "path";
import { fileURLToPath } from "url";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const LOCALE_DIR = path.resolve(__dirname, "../src/i18n/locales");

function getJsonFiles(dir) {
  return fs.readdirSync(dir).filter(file => file.endsWith(".json"));
}

function loadJson(filePath) {
  const content = fs.readFileSync(filePath, "utf-8");
  return JSON.parse(content);
}

function getDeepKeys(obj, prefix = "") {
  let keys = [];
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      const nextPrefix = prefix ? `${prefix}.${key}` : key;
      if (typeof obj[key] === "object" && obj[key] !== null && !Array.isArray(obj[key])) {
        keys = keys.concat(getDeepKeys(obj[key], nextPrefix));
      } else {
        keys.push({ key: nextPrefix, value: obj[key] });
      }
    }
  }
  return keys;
}

function verify() {
  const thDir = path.join(LOCALE_DIR, "th");
  const enDir = path.join(LOCALE_DIR, "en");

  if (!fs.existsSync(thDir) || !fs.existsSync(enDir)) {
    console.error("❌ Locale directories do not exist");
    process.exit(1);
  }

  const thFiles = getJsonFiles(thDir);
  const enFiles = getJsonFiles(enDir);

  // Check file lists match
  if (thFiles.length !== enFiles.length) {
    console.error(`❌ File count mismatch. TH: ${thFiles.length}, EN: ${enFiles.length}`);
    process.exit(1);
  }

  let hasError = false;

  for (const file of thFiles) {
    const thPath = path.join(thDir, file);
    const enPath = path.join(enDir, file);

    if (!fs.existsSync(enPath)) {
      console.error(`❌ English translation file missing: ${file}`);
      hasError = true;
      continue;
    }

    const thData = loadJson(thPath);
    const enData = loadJson(enPath);

    const thKeysObj = getDeepKeys(thData);
    const enKeysObj = getDeepKeys(enData);

    const thKeys = thKeysObj.map(k => k.key);
    const enKeys = enKeysObj.map(k => k.key);

    // 1. Check for missing keys in EN
    const missingInEn = thKeys.filter(k => !enKeys.includes(k));
    if (missingInEn.length > 0) {
      console.error(`❌ File ${file}: Keys missing in English:\n   ${missingInEn.join("\n   ")}`);
      hasError = true;
    }

    // 2. Check for missing keys in TH
    const missingInTh = enKeys.filter(k => !thKeys.includes(k));
    if (missingInTh.length > 0) {
      console.error(`❌ File ${file}: Keys missing in Thai:\n   ${missingInTh.join("\n   ")}`);
      hasError = true;
    }

    // 3. Check for empty values or invalid translations
    for (const item of thKeysObj) {
      if (item.value === "" || item.value === null || item.value === undefined) {
        console.error(`❌ File ${file}: Empty Thai translation for key: ${item.key}`);
        hasError = true;
      }
    }
    for (const item of enKeysObj) {
      if (item.value === "" || item.value === null || item.value === undefined) {
        console.error(`❌ File ${file}: Empty English translation for key: ${item.key}`);
        hasError = true;
      }
    }
  }

  if (hasError) {
    console.log("❌ Locale verification failed!");
    process.exit(1);
  } else {
    console.log("✅ All locales verified successfully! Trees are matching, values are non-empty, and keys are symmetric.");
  }
}

verify();

import * as App from "./wailsjs/go/main/App.js";
import { ClipboardSetText } from "./wailsjs/runtime/runtime.js";

const el = (id) => document.getElementById(id);

const UI = {
    ko: {
        appTitle: "시간 차이 계산기",
        appSubtitle: "한 입력란에 두 개의 날짜·시간을 넣으면 자동으로 인식합니다.",
        lblInput: "날짜·시간 입력",
        btnCalc: "계산",
        lblResults: "결과",
        lblSettings: "설정",
        lblFormats: "출력 형식",
        t_seconds: "초",
        t_milliseconds: "밀리초",
        t_mmss: "분:초",
        t_hhmmss: "시:분:초",
        t_ddhhmmss: "일 시:분:초",
        t_full: "전체(밀리초 포함)",
        t_custom: "사용자 정의 형식",
        lblCustom: "사용자 정의 형식 문자열",
        lblLanguage: "언어",
        copy: "복사",
        copied: "복사됨",
        copyFail: "복사 실패",
        langAuto: "자동 (OS)",
        lblZeroPadding: "앞자리 0 채우기 (시·분·초·밀리초)",
    },
    en: {
        appTitle: "Time Diff Calculator",
        appSubtitle: "Enter two date-times in one field; they are detected automatically.",
        lblInput: "Date-time input",
        btnCalc: "Calculate",
        lblResults: "Results",
        lblSettings: "Settings",
        lblFormats: "Output formats",
        t_seconds: "Seconds",
        t_milliseconds: "Milliseconds",
        t_mmss: "mm:ss",
        t_hhmmss: "hh:mm:ss",
        t_ddhhmmss: "dd hh:mm:ss",
        t_full: "Full (with ms)",
        t_custom: "Custom format",
        lblCustom: "Custom format string",
        lblLanguage: "Language",
        copy: "Copy",
        copied: "Copied",
        copyFail: "Copy failed",
        langAuto: "Auto (OS)",
        lblZeroPadding: "Zero-padding (hh, mm, ss, ms)",
    },
};

function applyChromeStrings(locale) {
    const L = UI[locale] || UI.en;
    const map = [
        ["appTitle", L.appTitle],
        ["appSubtitle", L.appSubtitle],
        ["lblInput", L.lblInput],
        ["btnCalc", L.btnCalc],
        ["lblResults", L.lblResults],
        ["lblSettings", L.lblSettings],
        ["lblFormats", L.lblFormats],
        ["t_seconds", L.t_seconds],
        ["t_milliseconds", L.t_milliseconds],
        ["t_mmss", L.t_mmss],
        ["t_hhmmss", L.t_hhmmss],
        ["t_ddhhmmss", L.t_ddhhmmss],
        ["t_full", L.t_full],
        ["t_custom", L.t_custom],
        ["lblCustom", L.lblCustom],
        ["lblLanguage", L.lblLanguage],
        ["lblZeroPadding", L.lblZeroPadding],
    ];
    for (const [id, text] of map) {
        const node = el(id);
        if (node) node.textContent = text;
    }
    const langSel = el("language");
    if (langSel && langSel.options[0]) {
        langSel.options[0].textContent = L.langAuto;
    }
}

function collectSettings() {
    const formats = {};
    document.querySelectorAll("input[data-fmt]").forEach((cb) => {
        formats[cb.getAttribute("data-fmt")] = cb.checked;
    });
    return {
        formats,
        customFormat: el("customFormat").value,
        language: el("language").value,
        zeroPadding: el("zeroPadding").checked,
    };
}

function applySettings(s) {
    const formats = s.formats || {};
    document.querySelectorAll("input[data-fmt]").forEach((cb) => {
        const k = cb.getAttribute("data-fmt");
        cb.checked = formats[k] !== false;
    });
    el("customFormat").value = s.customFormat || "";
    el("language").value = s.language || "auto";
    const zp = el("zeroPadding");
    if (zp) {
        zp.checked =
            s.zeroPadding === undefined || s.zeroPadding === null
                ? true
                : !!s.zeroPadding;
    }
}

let saveTimer;
function scheduleSave() {
    clearTimeout(saveTimer);
    saveTimer = setTimeout(() => {
        App.SaveSettings(collectSettings()).catch(console.error);
    }, 320);
}

function setValidation(v) {
    const node = el("validation");
    node.textContent = v.message || "";
    node.className = "validation";
    if (v.level === "warn") node.classList.add("warn");
    else if (v.level === "error") node.classList.add("error");
    else if (v.level === "ok") node.classList.add("ok");
}

async function runValidate() {
    const input = el("datetimeInput").value;
    const settings = collectSettings();
    try {
        const v = await App.ValidateInput(input, settings);
        setValidation(v);
        return v;
    } catch (e) {
        console.error(e);
        setValidation({ level: "error", message: String(e) });
        return { ready: false };
    }
}

function renderResults(resp) {
    const ul = el("results");
    ul.innerHTML = "";
    if (!resp.ok || !resp.results) return;
    const locale = resp.locale === "ko" ? "ko" : "en";
    const L = UI[locale] || UI.en;
    for (const line of resp.results) {
        const li = document.createElement("li");
        const left = document.createElement("div");
        const meta = document.createElement("div");
        meta.className = "result-meta";
        meta.textContent = line.label;
        const val = document.createElement("div");
        val.className = "result-value";
        val.textContent = line.value;
        left.appendChild(meta);
        left.appendChild(val);

        const btn = document.createElement("button");
        btn.type = "button";
        btn.className = "btn ghost";
        btn.textContent = L.copy;
        btn.addEventListener("click", async () => {
            try {
                await ClipboardSetText(line.value);
                btn.textContent = L.copied;
                setTimeout(() => (btn.textContent = L.copy), 1200);
            } catch {
                btn.textContent = L.copyFail;
                setTimeout(() => (btn.textContent = L.copy), 1200);
            }
        });

        li.appendChild(left);
        li.appendChild(btn);
        ul.appendChild(li);
    }
}

async function calculate() {
    const v = await runValidate();
    if (!v.ready) {
        el("results").innerHTML = "";
        return;
    }
    try {
        const resp = await App.Calculate(el("datetimeInput").value, collectSettings());
        setValidation({ level: resp.ok ? "ok" : "error", message: resp.error || resp.warning || v.message });
        renderResults(resp);
    } catch (e) {
        console.error(e);
        setValidation({ level: "error", message: String(e) });
    }
}

let validateTimer;
function debouncedValidate() {
    clearTimeout(validateTimer);
    validateTimer = setTimeout(async () => {
        const v = await runValidate();
        if (v.ready) {
            try {
                const resp = await App.Calculate(el("datetimeInput").value, collectSettings());
                renderResults(resp);
            } catch (e) {
                console.error(e);
            }
        } else {
            el("results").innerHTML = "";
        }
    }, 280);
}

async function init() {
    let sys = "en";
    try {
        sys = await App.SystemLocale();
    } catch (e) {
        console.warn(e);
    }
    let settings = {};
    try {
        settings = await App.LoadSettings();
    } catch (e) {
        console.warn(e);
    }
    applySettings(settings);
    const uiLocale =
        settings.language && settings.language !== "auto" ? settings.language : sys;
    applyChromeStrings(uiLocale === "ko" ? "ko" : "en");

    el("datetimeInput").addEventListener("input", () => {
        scheduleSave();
        debouncedValidate();
    });
    el("datetimeInput").addEventListener("keydown", (e) => {
        if (e.key === "Enter" && !e.shiftKey) {
            e.preventDefault();
            calculate();
        }
    });
    el("btnCalc").addEventListener("click", () => calculate());

    ["customFormat", "language"].forEach((id) => {
        el(id).addEventListener("change", () => {
            scheduleSave();
            debouncedValidate();
            if (id === "language") {
                const s = collectSettings();
                const loc =
                    s.language && s.language !== "auto" ? s.language : sys;
                applyChromeStrings(loc === "ko" ? "ko" : "en");
            }
        });
    });
    document.querySelectorAll("input[data-fmt]").forEach((cb) => {
        cb.addEventListener("change", () => {
            scheduleSave();
            debouncedValidate();
        });
    });
    el("zeroPadding").addEventListener("change", () => {
        scheduleSave();
        debouncedValidate();
    });

    el("datetimeInput").focus();
    debouncedValidate();
}

init().catch(console.error);

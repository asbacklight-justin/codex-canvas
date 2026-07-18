import './style.css';
import {
  ApplyTheme,
  GetStatus,
  ListPresets,
  RestoreOfficial,
  RevealDataFolder,
  SelectImage,
  SelectPreset,
} from '../wailsjs/go/main/App';

const translations = {
  zh: {
    brandSubtitle: '幕光', studio: '主题工坊', localFiles: '本地文件', language: '界面语言',
    followSystem: '跟随系统', chinese: '中文', english: 'English',
    safetyTitle: '可逆 · 本地优先', safetyLine1: '不修改 app.asar', safetyLine2: '不破坏应用签名',
    contact: '问题反馈与联系', unofficial: '非 OpenAI 官方产品 · v1.0.0',
    heroTitle: '为你的 Codex 换一束光', heroSubtitle: '选择图片，调好阅读层，然后安全启动。',
    detecting: '正在检测…', apply: '应用到 Codex', livePreview: 'LIVE PREVIEW', workspacePreview: '工作区预览',
    previewMode: '预览模式', themeRunning: '主题运行中', codexNotFound: '未找到 Codex', runtimeTheme: '运行时主题 · 不修改应用签名',
    newTask: '新任务', today: '今天', yesterday: '昨天', designPortfolio: '设计个人主页', fixUpload: '修复图片上传',
    projectDocs: '整理项目文档', themeResearch: '主题参数研究', localWorkspace: '本地工作区',
    samplePrompt: '为我设计一个克制、耐看的作品集主页。', sampleAnswer: '我会让内容保持清晰，把背景的氛围留在阅读层之后。',
    sampleCode: '精炼布局与视觉层级', continuePrompt: '继续描述你想要的效果…',
    chooseFirst: '先选择一张背景图片', wideRecommended: '推荐 1920×1080 以上的横向图片', chooseImage: '选择图片',
    privacy: '图片只保存在你的 Mac 上，不上传、不分析、不遥测。',
    backgroundImage: '背景图片', imageFormats: 'PNG · JPEG · WebP · GIF，最大 16 MB', chooseLocal: '选择本地图片',
    fileMeta: 'Finder 中挑选，文件不会上传', savedLocal: '已保存在本机', change: '更换',
    builtinTitle: '内置背景', builtinDesc: '原创生成 · 点击即可预览使用', loadingPresets: '正在加载预设…',
    readability: '阅读与层次', readabilityDesc: '先保证文字清楚，再保留图片氛围', overlay: '暗色遮罩',
    contentPanel: '内容面板', sidebarPanel: '侧栏面板', backgroundBlur: '背景模糊',
    details: '细节', detailsDesc: '少量调整通常更耐看', dialogRadius: '弹窗圆角', imageScale: '图片缩放',
    imagePosition: '图片位置', top: '顶部', center: '居中', bottom: '底部', accent: '强调色',
    restore: '恢复官方外观', closeApplyTitle: '关闭 Codex 并应用主题？',
    closeApplyBody: 'Codex Canvas 会正常退出当前 Codex，再用仅限本机的临时主题通道重新启动。未保存的输入内容可能丢失。',
    safeNote: '不会修改 Codex.app 或它的代码签名', cancel: '取消', closeApply: '关闭 Codex 并应用',
    processing: '处理中…', restarting: '正在重启 Codex…', restoring: '正在恢复…', selectImageFirst: '请先选择一张背景图片',
    appliedSuccess: '主题已应用，Codex 已使用新的背景和高对比文字。', restoredSuccess: '已恢复 Codex 官方外观。',
    presetLoadError: '内置背景加载失败', operationFailed: '操作失败', localOnly: '仅保存在本机',
    'preset.technology': '科技 · 霓虹地平线', 'preset.anime': '二次元 · 蓝调车站', 'preset.city': '城市 · 午夜天际线',
    'preset.nature': '自然 · 高山晨雾', 'preset.animal': '动物 · 森林赤狐',
  },
  en: {
    brandSubtitle: 'Luminous themes', studio: 'Theme Studio', localFiles: 'Local Files', language: 'Language',
    followSystem: 'Follow System', chinese: 'Chinese (中文)', english: 'English',
    safetyTitle: 'Reversible · Local-first', safetyLine1: 'Never modifies app.asar', safetyLine2: 'Preserves the app signature',
    contact: 'Support & feedback', unofficial: 'Independent product · Not affiliated with OpenAI · v1.0.0',
    heroTitle: 'Give your Codex a new light', heroSubtitle: 'Choose an image, tune readability, then launch safely.',
    detecting: 'Detecting…', apply: 'Apply to Codex', livePreview: 'LIVE PREVIEW', workspacePreview: 'Workspace preview',
    previewMode: 'Preview', themeRunning: 'Theme active', codexNotFound: 'Codex not found', runtimeTheme: 'Runtime theme · app signature preserved',
    newTask: 'New task', today: 'Today', yesterday: 'Yesterday', designPortfolio: 'Design portfolio', fixUpload: 'Fix image upload',
    projectDocs: 'Organize project docs', themeResearch: 'Research theme settings', localWorkspace: 'Local workspace',
    samplePrompt: 'Design a restrained, timeless portfolio for me.', sampleAnswer: 'I’ll keep the content clear and let the background atmosphere sit behind the reading layer.',
    sampleCode: 'Refine layout and visual hierarchy', continuePrompt: 'Describe the effect you want…',
    chooseFirst: 'Choose a background image first', wideRecommended: 'Landscape images at 1920×1080 or above work best', chooseImage: 'Choose image',
    privacy: 'Images stay on your Mac. No uploads, analysis, or telemetry.',
    backgroundImage: 'Background image', imageFormats: 'PNG · JPEG · WebP · GIF, up to 16 MB', chooseLocal: 'Choose a local image',
    fileMeta: 'Pick in Finder — nothing is uploaded', savedLocal: 'Saved locally', change: 'Change',
    builtinTitle: 'Built-in backgrounds', builtinDesc: 'Original artwork · click to preview', loadingPresets: 'Loading presets…',
    readability: 'Readability & depth', readabilityDesc: 'Keep text clear while preserving atmosphere', overlay: 'Dark overlay',
    contentPanel: 'Content panel', sidebarPanel: 'Sidebar panel', backgroundBlur: 'Background blur',
    details: 'Details', detailsDesc: 'Small adjustments usually age better', dialogRadius: 'Dialog radius', imageScale: 'Image scale',
    imagePosition: 'Image position', top: 'Top', center: 'Center', bottom: 'Bottom', accent: 'Accent',
    restore: 'Restore official', closeApplyTitle: 'Quit Codex and apply this theme?',
    closeApplyBody: 'Codex Canvas will quit Codex normally and relaunch it through a temporary local theme channel. Unsent input may be lost.',
    safeNote: 'Codex.app and its code signature are never modified', cancel: 'Cancel', closeApply: 'Quit Codex and apply',
    processing: 'Working…', restarting: 'Restarting Codex…', restoring: 'Restoring…', selectImageFirst: 'Choose a background image first',
    appliedSuccess: 'Theme applied with the new background and high-contrast text.', restoredSuccess: 'The official Codex appearance has been restored.',
    presetLoadError: 'Could not load built-in backgrounds', operationFailed: 'Operation failed', localOnly: 'Stored only on this Mac',
    'preset.technology': 'Tech · Neon Horizon', 'preset.anime': 'Anime · Blue-hour Station', 'preset.city': 'City · Midnight Skyline',
    'preset.nature': 'Nature · Alpine Mist', 'preset.animal': 'Animal · Forest Fox',
  },
};

const languageKey = 'codexCanvasLanguage';
const storedPreference = localStorage.getItem(languageKey);
let languagePreference = ['system', 'zh', 'en'].includes(storedPreference) ? storedPreference : 'system';

function detectSystemLanguage() {
  const languages = Array.isArray(navigator.languages) && navigator.languages.length
    ? navigator.languages
    : (navigator.language ? [navigator.language] : []);
  if (!languages.length) return 'en';
  return String(languages[0]).toLowerCase().startsWith('zh') ? 'zh' : 'en';
}

const currentLanguage = languagePreference === 'system' ? detectSystemLanguage() : languagePreference;
const t = (key) => translations[currentLanguage][key] || translations.en[key] || key;
document.documentElement.lang = currentLanguage === 'zh' ? 'zh-CN' : 'en';

const defaultTheme = {
  imagePath: '', imageName: '', overlay: 48, surfaceOpacity: 72, sidebarOpacity: 82,
  blur: 0, radius: 18, scale: 100, position: 'center', accent: '#8b7cff', active: false,
};

let theme = { ...defaultTheme };
let previewUrl = '';
let status = null;
let busy = false;
let presets = [];
let selectedPresetID = '';
let imageMetaText = t('fileMeta');

document.querySelector('#app').innerHTML = `
  <div class="shell">
    <aside class="sidebar">
      <div class="brand"><div class="brand-mark"><span></span><span></span></div><div><strong>Codex Canvas</strong><small>${t('brandSubtitle')}</small></div></div>
      <nav>
        <button class="nav-item active" type="button"><span class="nav-icon">✦</span><span>${t('studio')}</span></button>
        <button class="nav-item" id="data-folder" type="button"><span class="nav-icon">⌁</span><span>${t('localFiles')}</span></button>
      </nav>
      <div class="language-card">
        <label for="language-select">${t('language')}</label>
        <select id="language-select" aria-label="${t('language')}">
          <option value="system">${t('followSystem')}</option>
          <option value="zh">${t('chinese')}</option>
          <option value="en">${t('english')}</option>
        </select>
      </div>
      <section class="sidebar-presets" aria-labelledby="sidebar-presets-title">
        <div class="sidebar-section-title">
          <strong id="sidebar-presets-title">${t('builtinTitle')}</strong>
          <small>${t('builtinDesc')}</small>
        </div>
        <div class="preset-grid" id="preset-grid"><div class="preset-loading">${t('loadingPresets')}</div></div>
      </section>
      <div class="side-spacer"></div>
      <div class="safety-card"><div class="safety-icon">✓</div><div><strong>${t('safetyTitle')}</strong><p>${t('safetyLine1')}<br>${t('safetyLine2')}</p></div></div>
      <div class="contact-card"><span>${t('contact')}</span><strong>asbacklight@gmail.com</strong></div>
      <div class="unofficial">${t('unofficial')}</div>
    </aside>

    <main class="workspace">
      <header class="topbar">
        <div><h1>${t('heroTitle')}</h1><p>${t('heroSubtitle')}</p></div>
        <div class="top-actions">
          <div class="codex-status" id="codex-status"><span class="status-dot"></span><div><strong>${t('detecting')}</strong><small>Codex Desktop</small></div></div>
          <button class="primary top-apply" id="apply-top" type="button"><span>✦</span> ${t('apply')}</button>
        </div>
      </header>

      <section class="content-grid">
        <div class="preview-column">
          <div class="section-heading"><div><span class="eyebrow">${t('livePreview')}</span><h2>${t('workspacePreview')}</h2></div><div class="theme-state" id="theme-state">${t('previewMode')}</div></div>
          <div class="codex-preview" id="preview">
            <div class="preview-image" id="preview-image"></div><div class="preview-dim" id="preview-dim"></div>
            <div class="mock-sidebar" id="mock-sidebar">
              <div class="mock-logo"><i></i><span>Codex</span></div><button>＋&nbsp;&nbsp;${t('newTask')}</button><small>${t('today')}</small>
              <p class="selected">${t('designPortfolio')}</p><p>${t('fixUpload')}</p><p>${t('projectDocs')}</p><small>${t('yesterday')}</small><p>${t('themeResearch')}</p>
              <div class="mock-user"><span>C</span><div>Canvas User<small>${t('localWorkspace')}</small></div></div>
            </div>
            <div class="mock-main" id="mock-main">
              <div class="mock-toolbar"><span>${t('designPortfolio')}</span><b>•••</b></div>
              <div class="mock-conversation"><div class="mock-prompt">${t('samplePrompt')}</div><div class="mock-answer"><span class="spark">✦</span><p>${t('sampleAnswer')}</p></div><div class="mock-code"><div><i></i><i></i><i></i><span>portfolio.tsx</span></div><p><em>+</em> 128&nbsp;&nbsp; ${t('sampleCode')}</p></div></div>
              <div class="mock-composer" id="mock-composer"><span>${t('continuePrompt')}</span><button>↑</button></div>
            </div>
            <div class="empty-preview" id="empty-preview"><div class="empty-orbit"><span>✦</span></div><strong>${t('chooseFirst')}</strong><p>${t('wideRecommended')}</p><button class="choose-inline" type="button">${t('chooseImage')}</button></div>
          </div>
          <div class="privacy-strip"><span>▣</span> ${t('privacy')}</div>
        </div>

        <aside class="controls-panel">
          <div class="controls-scroll">
            <section class="control-section image-control">
              <div class="control-title"><span>01</span><div><strong>${t('backgroundImage')}</strong><small>${t('imageFormats')}</small></div></div>
              <button class="file-card" id="choose-image" type="button"><div class="file-thumb" id="file-thumb">⌁</div><div><strong id="file-name">${t('chooseLocal')}</strong><small id="file-meta">${t('fileMeta')}</small></div><span>${t('change')}</span></button>
            </section>
            <section class="control-section">
              <div class="control-title"><span>02</span><div><strong>${t('readability')}</strong><small>${t('readabilityDesc')}</small></div></div>
              ${rangeControl('overlay', t('overlay'), 0, 90, 1, '%')}${rangeControl('surfaceOpacity', t('contentPanel'), 15, 100, 1, '%')}${rangeControl('sidebarOpacity', t('sidebarPanel'), 15, 100, 1, '%')}${rangeControl('blur', t('backgroundBlur'), 0, 30, 1, ' px')}
            </section>
            <section class="control-section compact-section">
              <div class="control-title"><span>03</span><div><strong>${t('details')}</strong><small>${t('detailsDesc')}</small></div></div>
              ${rangeControl('radius', t('dialogRadius'), 0, 32, 1, ' px')}${rangeControl('scale', t('imageScale'), 100, 150, 1, '%')}
              <div class="field-row"><label>${t('imagePosition')}</label><div class="segmented" id="position-control"><button data-position="top" type="button">${t('top')}</button><button data-position="center" type="button" class="active">${t('center')}</button><button data-position="bottom" type="button">${t('bottom')}</button></div></div>
              <div class="field-row color-row"><label for="accent">${t('accent')}</label><div class="color-picker"><input id="accent" type="color" value="#8b7cff"><span id="accent-value">#8B7CFF</span></div></div>
            </section>
          </div>
          <div class="action-bar"><button class="secondary" id="restore" type="button">${t('restore')}</button><button class="primary" id="apply" type="button"><span>✦</span> ${t('apply')}</button></div>
        </aside>
      </section>
    </main>
  </div>
  <div class="modal-backdrop" id="confirm-modal" hidden><div class="modal"><div class="modal-icon">✦</div><h3>${t('closeApplyTitle')}</h3><p>${t('closeApplyBody')}</p><div class="modal-note"><span>✓</span> ${t('safeNote')}</div><div class="modal-actions"><button class="secondary" id="cancel-apply">${t('cancel')}</button><button class="primary" id="confirm-apply">${t('closeApply')}</button></div></div></div>
  <div class="toast" id="toast" role="status"></div>
`;

function rangeControl(id, label, min, max, step, suffix) {
  return `<div class="range-control"><div><label for="${id}">${label}</label><output id="${id}-value">—</output></div><input id="${id}" type="range" min="${min}" max="${max}" step="${step}" data-suffix="${suffix}"></div>`;
}

const $ = (selector) => document.querySelector(selector);
$('#language-select').value = languagePreference;

function setBusy(value, label = '') {
  busy = value;
  ['#apply', '#apply-top', '#restore', '#choose-image'].forEach((selector) => { $(selector).disabled = value; });
  const content = value ? `<span class="spinner"></span>${label || t('processing')}` : `<span>✦</span> ${t('apply')}`;
  $('#apply').innerHTML = content;
  $('#apply-top').innerHTML = content;
}

function translatedImageName(name) {
  for (const id of ['technology', 'anime', 'city', 'nature', 'animal']) {
    if (name === id || name === translations.zh[`preset.${id}`] || name === translations.en[`preset.${id}`]) return t(`preset.${id}`);
  }
  return name;
}

function render() {
  const hasImage = Boolean(previewUrl);
  $('#preview').classList.toggle('has-image', hasImage);
  $('#preview-image').style.backgroundImage = hasImage ? `url("${previewUrl}")` : '';
  $('#preview-image').style.backgroundPosition = theme.position;
  $('#preview-image').style.filter = `blur(${theme.blur / 3}px)`;
  $('#preview-image').style.transform = `scale(${theme.scale / 100})`;
  $('#preview-dim').style.background = `rgba(4, 5, 10, ${theme.overlay / 100})`;
  $('#mock-main').style.background = `rgba(11, 12, 19, ${theme.surfaceOpacity / 100})`;
  $('#mock-sidebar').style.background = `rgba(7, 8, 13, ${theme.sidebarOpacity / 100})`;
  $('#mock-composer').style.borderRadius = `${Math.min(theme.radius, 22)}px`;
  $('#preview').style.setProperty('--accent', theme.accent);
  $('#file-thumb').style.backgroundImage = hasImage ? `url("${previewUrl}")` : '';
  $('#file-thumb').classList.toggle('has-image', hasImage);
  $('#file-thumb').textContent = hasImage ? '' : '⌁';
  $('#file-name').textContent = translatedImageName(theme.imageName) || t('chooseLocal');
  $('#file-meta').textContent = imageMetaText;
  ['overlay', 'surfaceOpacity', 'sidebarOpacity', 'blur', 'radius', 'scale'].forEach((id) => {
    const input = $(`#${id}`);
    input.value = theme[id];
    $(`#${id}-value`).textContent = `${theme[id]}${input.dataset.suffix}`;
    input.style.setProperty('--range-progress', `${((theme[id] - Number(input.min)) / (Number(input.max) - Number(input.min))) * 100}%`);
  });
  $('#accent').value = theme.accent;
  $('#accent-value').textContent = theme.accent.toUpperCase();
  document.querySelectorAll('[data-position]').forEach((button) => button.classList.toggle('active', button.dataset.position === theme.position));
  document.querySelectorAll('[data-preset]').forEach((button) => button.classList.toggle('active', button.dataset.preset === selectedPresetID));
  $('#theme-state').textContent = theme.active ? t('themeRunning') : t('previewMode');
  $('#theme-state').classList.toggle('active', theme.active);
}

function renderStatus() {
  if (!status) return;
  const card = $('#codex-status');
  card.classList.toggle('found', status.codexFound);
  card.classList.toggle('active', status.active);
  card.querySelector('strong').textContent = status.codexFound ? (status.active ? t('themeRunning') : `Codex ${status.codexVersion || ''}`.trim()) : t('codexNotFound');
  card.querySelector('small').textContent = status.codexFound ? t('runtimeTheme') : t('codexNotFound');
  $('#apply').disabled = !status.codexFound;
  $('#apply-top').disabled = !status.codexFound;
}

function renderPresets() {
  const grid = $('#preset-grid');
  if (!presets.length) return;
  grid.innerHTML = presets.map((preset) => `<button class="preset-card" type="button" data-preset="${preset.id}"><span class="preset-image" style="background-image:url('${preset.previewUrl}')"></span><strong>${t(`preset.${preset.id}`)}</strong><i>✓</i></button>`).join('');
}

async function loadStatus() {
  try {
    status = await GetStatus();
    theme = { ...defaultTheme, ...(status.savedTheme || {}) };
    previewUrl = status.previewUrl || '';
    imageMetaText = theme.imagePath ? t('savedLocal') : t('fileMeta');
    renderStatus(); render();
  } catch (error) { showToast(cleanError(error), true); }
}

async function loadPresets() {
  try { presets = await ListPresets(); renderPresets(); render(); }
  catch (error) { $('#preset-grid').innerHTML = `<div class="preset-loading error">${t('presetLoadError')}</div>`; }
}

function useSelection(selection, displayName, presetID = '') {
  theme.imagePath = selection.path;
  theme.imageName = displayName || selection.name;
  previewUrl = selection.previewUrl;
  selectedPresetID = presetID;
  imageMetaText = `${formatBytes(selection.size)} · ${t('localOnly')}`;
  render();
}

async function chooseImage() {
  if (busy) return;
  try { const selection = await SelectImage(currentLanguage); if (!selection.cancelled) useSelection(selection, selection.name); }
  catch (error) { showToast(cleanError(error), true); }
}

async function choosePreset(id) {
  if (busy) return;
  try { const selection = await SelectPreset(id); useSelection(selection, t(`preset.${id}`), id); }
  catch (error) { showToast(cleanError(error), true); }
}

function applyTheme() {
  if (!theme.imagePath) { showToast(t('selectImageFirst'), true); return; }
  $('#confirm-modal').hidden = false;
}

async function confirmApply() {
  $('#confirm-modal').hidden = true; setBusy(true, t('restarting'));
  try { await ApplyTheme({ ...theme }); theme.active = true; if (status) status.active = true; render(); renderStatus(); showToast(t('appliedSuccess')); }
  catch (error) { showToast(cleanError(error), true); }
  finally { setBusy(false); renderStatus(); }
}

async function restoreOfficial() {
  if (busy) return; setBusy(true, t('restoring'));
  try { await RestoreOfficial(); theme.active = false; if (status) status.active = false; render(); renderStatus(); showToast(t('restoredSuccess')); }
  catch (error) { showToast(cleanError(error), true); }
  finally { setBusy(false); renderStatus(); }
}

let toastTimer;
function showToast(message, isError = false) {
  const toast = $('#toast'); toast.textContent = message; toast.classList.toggle('error', isError); toast.classList.add('show');
  clearTimeout(toastTimer); toastTimer = setTimeout(() => toast.classList.remove('show'), 4200);
}

function cleanError(error) {
  const raw = String(error?.message || error || t('operationFailed')).replace(/^Error:\s*/i, '');
  if (currentLanguage === 'zh') return raw;
  const mappings = [
    ['请先选择一张背景图片', 'Choose a background image first'], ['图片不能超过 16 MB', 'The image must be 16 MB or smaller'],
    ['未找到 Codex Desktop', 'Codex Desktop was not found'], ['所选图片已不存在', 'The selected image no longer exists'],
    ['图片格式不受支持', 'This image format is not supported'], ['Codex 仍在运行', 'Codex is still running; quit it manually and try again'],
    ['无法启动 Codex', 'Could not launch Codex'], ['主题应用失败', 'Could not apply the theme'], ['操作失败', 'Operation failed'],
  ];
  const match = mappings.find(([source]) => raw.includes(source));
  return match ? match[1] : raw;
}

function formatBytes(bytes) {
  if (!bytes) return '0 KB';
  return bytes > 1024 * 1024 ? `${(bytes / 1024 / 1024).toFixed(1)} MB` : `${Math.ceil(bytes / 1024)} KB`;
}

['overlay', 'surfaceOpacity', 'sidebarOpacity', 'blur', 'radius', 'scale'].forEach((id) => $(`#${id}`).addEventListener('input', (event) => { theme[id] = Number(event.target.value); render(); }));
$('#accent').addEventListener('input', (event) => { theme.accent = event.target.value; render(); });
$('#position-control').addEventListener('click', (event) => { const button = event.target.closest('[data-position]'); if (button) { theme.position = button.dataset.position; render(); } });
$('#preset-grid').addEventListener('click', (event) => { const button = event.target.closest('[data-preset]'); if (button) choosePreset(button.dataset.preset); });
$('#choose-image').addEventListener('click', chooseImage); $('.choose-inline').addEventListener('click', chooseImage);
$('#apply').addEventListener('click', applyTheme); $('#apply-top').addEventListener('click', applyTheme); $('#restore').addEventListener('click', restoreOfficial);
$('#cancel-apply').addEventListener('click', () => { $('#confirm-modal').hidden = true; }); $('#confirm-apply').addEventListener('click', confirmApply);
$('#confirm-modal').addEventListener('click', (event) => { if (event.target === event.currentTarget) event.currentTarget.hidden = true; });
$('#data-folder').addEventListener('click', async () => { try { await RevealDataFolder(); } catch (error) { showToast(cleanError(error), true); } });
$('#language-select').addEventListener('change', (event) => { localStorage.setItem(languageKey, event.target.value); window.location.reload(); });

loadStatus();
loadPresets();

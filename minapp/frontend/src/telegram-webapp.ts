// Telegram WebApp API types and wrapper
interface TelegramWebApp {
  initData: string;
  initDataUnsafe: any;
  version: string;
  platform: string;
  colorScheme: 'light' | 'dark';
  themeParams: any;
  isExpanded: boolean;
  viewportHeight: number;
  viewportStableHeight: number;
  headerColor: string;
  backgroundColor: string;
  isClosingConfirmationEnabled: boolean;
  BackButton: any;
  MainButton: any;
  HapticFeedback: any;
  
  ready(): void;
  expand(): void;
  close(): void;
  setHeaderColor(color: string): void;
  setBackgroundColor(color: string): void;
  enableClosingConfirmation(): void;
  disableClosingConfirmation(): void;
  showAlert(message: string): void;
  showConfirm(message: string): void;
  showPopup(params: any): void;
  showScanQrPopup(params: any): void;
  closeScanQrPopup(): void;
  readTextFromClipboard(): void;
  requestWriteAccess(): void;
  requestContact(): void;
  openLink(url: string, options?: { try_instant_view?: boolean }): void;
  openTelegramLink(url: string): void;
  openInvoice(url: string): void;
}

declare global {
  interface Window {
    Telegram?: {
      WebApp: TelegramWebApp;
    };
  }
}

export const WebApp: TelegramWebApp = window.Telegram?.WebApp || {
  initData: '',
  initDataUnsafe: {},
  version: '6.0',
  platform: 'unknown',
  colorScheme: 'light',
  themeParams: {},
  isExpanded: false,
  viewportHeight: window.innerHeight,
  viewportStableHeight: window.innerHeight,
  headerColor: '#ffffff',
  backgroundColor: '#ffffff',
  isClosingConfirmationEnabled: false,
  BackButton: {},
  MainButton: {},
  HapticFeedback: {},
  
  ready: () => console.log('WebApp ready'),
  expand: () => console.log('WebApp expand'),
  close: () => console.log('WebApp close'),
  setHeaderColor: (color: string) => console.log('setHeaderColor', color),
  setBackgroundColor: (color: string) => console.log('setBackgroundColor', color),
  enableClosingConfirmation: () => console.log('enableClosingConfirmation'),
  disableClosingConfirmation: () => console.log('disableClosingConfirmation'),
  showAlert: (message: string) => alert(message),
  showConfirm: (message: string) => window.confirm(message),
  showPopup: (params: any) => console.log('showPopup', params),
  showScanQrPopup: (params: any) => console.log('showScanQrPopup', params),
  closeScanQrPopup: () => console.log('closeScanQrPopup'),
  readTextFromClipboard: () => console.log('readTextFromClipboard'),
  requestWriteAccess: () => console.log('requestWriteAccess'),
  requestContact: () => console.log('requestContact'),
  openLink: (url: string) => window.open(url, '_blank'),
  openTelegramLink: (url: string) => console.log('openTelegramLink', url),
  openInvoice: (url: string) => console.log('openInvoice', url),
};

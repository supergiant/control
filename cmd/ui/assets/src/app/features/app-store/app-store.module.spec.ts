import { AppStoreModule } from './app-store.module';

describe('AppStoreModule', () => {
  let appStoreModule: AppStoreModule;

  beforeEach(() => {
    appStoreModule = new AppStoreModule();
  });

  it('should create an instance', () => {
    expect(appStoreModule).toBeTruthy();
  });
});

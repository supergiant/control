import { SguiPage } from './app.po';

describe('sgui App', () => {
  let page: SguiPage;

  beforeEach(() => {
    page = new SguiPage();
  });

  it('should display welcome message', () => {
    page.navigateTo();
    expect(page.getParagraphText()).toEqual('Welcome to app!');
  });
});

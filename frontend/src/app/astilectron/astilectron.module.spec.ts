import { AstilectronModule } from './astilectron.module';

describe('AstilectronModule', () => {
  let astilectronModule: AstilectronModule;

  beforeEach(() => {
    astilectronModule = new AstilectronModule();
  });

  it('should create an instance', () => {
    expect(astilectronModule).toBeTruthy();
  });
});

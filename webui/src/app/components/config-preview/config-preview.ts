import { Component, computed, input, ChangeDetectionStrategy } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatIconModule } from '@angular/material/icon';
import { MatTooltipModule } from '@angular/material/tooltip';
import * as yaml from 'js-yaml';
import { Config } from '../../models/config';

@Component({
  selector: 'app-config-preview',
  imports: [MatIconModule, MatButtonModule, MatTooltipModule],
  templateUrl: './config-preview.html',
  changeDetection: ChangeDetectionStrategy.Eager,
  styleUrl: './config-preview.css',
})
export class ConfigPreviewComponent {
  config = input<Config | undefined>();

  yamlString = computed(() => {
    if (this.config() === undefined) {
      return '';
    }
    return yaml.dump(this.config());
  });

  downloadConfig() {
    const config = this.config();
    if (config) {
      this.downloadYAML(config, 'config.yaml');
    }
  }

  downloadJSON(data: object, filename: string) {
    const content = JSON.stringify(data, null, 2);
    const dataUri = URL.createObjectURL(
      new Blob([content], {
        type: 'application/json;charset=utf-8',
      }),
    );
    const dummyLink = document.createElement('a');
    dummyLink.href = dataUri;
    dummyLink.download = filename;

    document.body.appendChild(dummyLink);
    dummyLink.click();
    document.body.removeChild(dummyLink);
  }

  downloadYAML(data: object, filename: string) {
    const content = yaml.dump(data);
    const dataUri = URL.createObjectURL(
      new Blob([content], {
        type: 'application/yaml;charset=utf-8',
      }),
    );
    const dummyLink = document.createElement('a');
    dummyLink.href = dataUri;
    dummyLink.download = filename;

    document.body.appendChild(dummyLink);
    dummyLink.click();
    document.body.removeChild(dummyLink);
  }
}

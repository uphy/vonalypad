<template>
  <div v-if="recipe">
    <v-card>
      <v-img
        class="white--text align-end"
        height="200px"
        :src="recipe.image"
      >
        <v-card-title>{{ recipe.title}}</v-card-title>
      </v-img>

      <v-card-text>
        {{ recipe.description }}
      </v-card-text>

      <v-divider></v-divider>
      <div style="padding: 1rem">
        <v-checkbox
          v-for="(ingredient, i) in recipe.ingredients"
          :key="i"
          :label="`${ingredient.name}: ${ingredient.quantity}`"
          dense
        ></v-checkbox>
      </div>

      <v-divider></v-divider>
      <v-tabs v-model="tab">
        <v-tabs-slider></v-tabs-slider>
        <v-tab>Overview</v-tab>
        <v-tab>Steps</v-tab>
      </v-tabs>

      <div style="padding: 1rem">
        <v-tabs-items v-model="tab">
          <v-tab-item>
            <ol>
              <li
                v-for="(step, i) in recipe.steps"
                :key="i"
                class="mb-8"
              >{{ step.step}}
                <v-img
                  v-if="step.image"
                  :src="step.image"
                  max-width="300px"
                />
              </li>
            </ol>
          </v-tab-item>
          <v-tab-item>
            <div class="text-center">
              <v-pagination
                v-model="step"
                :length="recipe.steps.length"
              ></v-pagination>
            </div>
            <div>{{ currentStep.step }}</div>
            <v-img
              v-if="currentStep.image"
              :src="currentStep.image"
              max-width="300px"
            />
          </v-tab-item>
        </v-tabs-items>
      </div>
    </v-card>
  </div>
</template>

<script>
import axios from "axios";
import NoSleep from "nosleep.js";

const noSleep = new NoSleep();

class VoiceControl {
  constructor(onVoice) {
    const SpeechRecognition =
      window.webkitSpeechRecognition || window.SpeechRecognition;
    if (SpeechRecognition === null || SpeechRecognition === undefined) {
      this.recognition = null;
      return;
    }
    this.recognition = new SpeechRecognition();
    this.recognition.continuous = true;
    this.recognition.lang = "ja-JP";

    this.recognition.onresult = (event) => {
      const results = event.results;
      onVoice(results[results.length - 1][0].transcript);
    };
  }

  start() {
    if (this.recognition === null) {
      return;
    }
    this.recognition.start();
  }

  stop() {
    if (this.recognition === null) {
      return;
    }
    this.recognition.stop();
  }
}

export default {
  data: () => {
    return {
      recipe: null,
      dialog: false,
      tab: 0,
      step: 1,
      currentStep: null,
      voiceControl: null,
    };
  },
  async mounted() {
    this.voiceControl = new VoiceControl((transcript) => {
      switch (transcript) {
        case "次":
          this.nextStep();
          break;
        case "前":
          this.prevStep();
          break;
        case "最初":
          this.step = 1;
          break;
        case "最後":
          this.step = this.recipe.steps.length;
          break;
      }
    });
    const recipeId = this.$route.params["recipeId"];
    const resp = await axios.get(`./api/recipes/${recipeId}`);
    this.recipe = resp.data;
    this.currentStep =
      this.recipe.steps.length > 0 ? this.recipe.steps[0] : null;
    this.dialog = true;
    this.step = 1;
    noSleep.enable();
    this.voiceControl.start();
  },
  destroyed() {
    this.voiceControl.stop();
    noSleep.disable();
    this.dialog = false;
    this.recipe = null;
  },
  watch: {
    step() {
      this.currentStep = this.recipe.steps[this.step - 1];
    },
  },
  methods: {
    nextStep() {
      this.step = Math.min(this.step + 1, this.recipe.steps.length);
    },
    prevStep() {
      this.step = Math.max(this.step - 1, 1);
    },
  },
};
</script>
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
        <v-tab @click="setStep(1)">Steps</v-tab>
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
    <div v-if="recognitionTexts.length > 0">
      <ul>
        <li
          v-for="(t,i) in recognitionTexts"
          :key="i"
        >{{ t }}</li>
      </ul>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import NoSleep from "nosleep.js";

const noSleep = new NoSleep();

class VoiceControl {
  constructor(onVoice, print) {
    this.stopping = true;

    const SpeechRecognition =
      window.webkitSpeechRecognition || window.SpeechRecognition;
    if (SpeechRecognition === null || SpeechRecognition === undefined) {
      this.recognition = null;
      return;
    }
    this.recognition = new SpeechRecognition();
    this.recognition.lang = "ja-JP";

    this.recognition.onresult = (event) => {
      print("onResult");
      const results = event.results;
      if (results.length > 0) {
        const result = results[results.length - 1];
        if (result.length > 0) {
          try {
            onVoice(result[0].transcript);
          } catch (e) {
            console.log(e);
          }
        }
      }
    };
    this.recognition.onend = () => {
      print("onEnd");
      if (!this.stopping) {
        print("restarting");
        this.start();
      }
    };
    this.recognition.onaudioend = () => {
      print("onAudioEnd");
    };
    this.recognition.onnomatch = () => {
      print("onNoMatch");
    };
    this.recognition.onsoundstart = () => {
      print("onSoundStart");
    };
    this.recognition.onsoundend = () => {
      print("onSoundEnd");
    };
    this.recognition.onspeechstart = () => {
      print("onSpeechStart");
    };
    this.recognition.onspeechend = () => {
      print("onSpeechEnd");
    };
  }

  start() {
    if (this.recognition === null) {
      return;
    }
    this.stopping = false;
    this.recognition.start();
  }

  stop() {
    if (this.recognition === null) {
      return;
    }
    this.stopping = true;
    this.recognition.stop();
  }
}

class TextToSpeech {
  constructor() {
    this.speechSynthesis = window.speechSynthesis;
  }

  speak(text) {
    if (this.speechSynthesis === null || this.speechSynthesis === undefined) {
      return;
    }
    this.speechSynthesis.cancel();
    const uttr = new SpeechSynthesisUtterance(text);
    this.speechSynthesis.speak(uttr);
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
      textToSpeech: null,
      debug: false,
      recognitionTexts: [],
    };
  },
  async mounted() {
    if (this.$route.query["debug"] === "true") {
      this.debug = true;
    }
    this.voiceControl = new VoiceControl(
      (transcript) => {
        transcript = transcript.toLowerCase();
        if (this.debug) {
          this.recognitionTexts.push(transcript);
        }
        if (transcript.includes("次") || transcript.includes("next")) {
          this.step = Math.min(this.step + 1, this.recipe.steps.length);
        } else if (
          transcript.includes("前") ||
          transcript.includes("プレビアス") ||
          transcript.includes("previous")
        ) {
          this.step = Math.max(this.step - 1, 1);
        } else if (
          transcript.includes("最初") ||
          transcript.includes("first") ||
          transcript.includes("ファースト")
        ) {
          this.step = 1;
        } else if (
          transcript.includes("最後") ||
          transcript.includes("last") ||
          transcript.includes("ラスト")
        ) {
          this.step = this.recipe.steps.length;
        }
      },
      (text) => {
        if (this.debug) {
          this.recognitionTexts.push(text);
        }
      }
    );
    this.textToSpeech = new TextToSpeech();
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
      this.setStep(this.step);
    },
  },
  methods: {
    setStep(step) {
      this.currentStep = this.recipe.steps[step - 1];
      this.textToSpeech.speak(this.currentStep.step);
      this.step = step;
    },
  },
};
</script>
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
        <div>{{ recipe.description }}</div>
        <v-btn
          color="#ff9933"
          class="ma-2 white--text"
          small
          @click="openInCookpad"
        >
          Cookpad
          <v-icon
            right
            dark
          >
            mdi-open-in-new
          </v-icon>
        </v-btn>
        <tags :recipeId="recipe.id" />
      </v-card-text>

      <v-divider></v-divider>
      <div style="padding: 1rem">
        <span>分量: {{ recipe.servingsFor }}</span>
        <div id="ingredients">
          <v-checkbox
            v-for="(ingredient, i) in recipe.ingredients"
            :key="i"
            :label="`${ingredient.name}: ${ingredient.quantity}`"
            dense
          ></v-checkbox>
        </div>
      </div>

      <v-divider></v-divider>

      <div style="padding: 1rem">
        <ol>
          <li
            v-for="(s, i) in recipe.steps"
            :key="i"
            class="mb-8"
            :id="`step-${i+1}`"
          >
            <span :style="{'background-color':step-1 === i ? '#ffa':null}">{{ s.step}}</span>
            <v-img
              v-if="s.image"
              :src="s.image"
              max-width="300px"
            />
          </li>
        </ol>
      </div>
    </v-card>
    <div style="margin-top: 4rem;"></div>
    <v-btn
      color="green"
      fab
      large
      dark
      bottom
      right
      fixed
      @click="nextStep"
    >
      <v-icon>mdi-arrow-right-bold</v-icon>
    </v-btn>
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
import Tags from "@/components/Tags.vue";

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
    this.enable = true;
  }

  speak(text) {
    if (this.speechSynthesis === null || this.speechSynthesis === undefined) {
      return;
    }
    if (!this.enable) {
      return;
    }
    this.speechSynthesis.cancel();
    const uttr = new SpeechSynthesisUtterance(text);
    this.speechSynthesis.speak(uttr);
  }
}

export default {
  components: {
    Tags,
  },
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
      tags: [],
      allTags: [],
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
          this.nextStep();
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
        } else if (
          transcript.includes("材料") ||
          transcript.includes("調味料")
        ) {
          this.focusElement("ingredients");
        } else if (
          transcript.includes("戻る") ||
          transcript.includes("手順") ||
          transcript.includes("現在") ||
          transcript.includes("今")
        ) {
          this.focusStep(this.step);
        }
      },
      (text) => {
        if (this.debug) {
          this.recognitionTexts.push(text);
        }
      }
    );
    this.textToSpeech = new TextToSpeech();
    this.textToSpeech.enable = false;
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
    nextStep() {
      this.step = Math.min(this.step + 1, this.recipe.steps.length);
    },
    setStep(step) {
      this.currentStep = this.recipe.steps[step - 1];
      this.textToSpeech.speak(this.currentStep.step);
      this.step = step;
      this.focusStep(step);
    },
    focusStep(step) {
      this.focusElement(`step-${step}`);
    },
    focusElement(elm) {
      document.getElementById(elm).scrollIntoView({
        behavior: "smooth",
        block: "center",
        inline: "start",
      });
    },
    openInCookpad() {
      window.open(`https://cookpad.com/recipe/${this.recipe.id}`, "_blank");
    },
  },
};
</script>
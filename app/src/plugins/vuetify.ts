import Vue from 'vue'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import { VuetifyPreset } from 'vuetify/types/presets'

Vue.use(Vuetify)

const opt: VuetifyPreset = {
    icons: {
        iconfont: 'md'
    },
    theme: {
        dark: true
    }
}

export default new Vuetify(opt)

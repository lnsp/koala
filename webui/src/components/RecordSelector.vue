<template>
  <Listbox :modelValue="modelValue"
           @update:modelValue="$emit('update:modelValue', $event)"
           as="div"
           class="relative">
    <ListboxButton class="relative w-full bg-white h-10 border border-gray-300 rounded-md shadow-sm text-left cursor-default focus:outline-none focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 text-sm">
      <span class="flex items-center">
        <span class="ml-3 block truncate">{{ modelValue }}</span>
      </span>
      <span class="ml-3 absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
        <SelectorIcon class="h-5 w-5 text-gray-400"
                      aria-hidden="true" />
      </span>
    </ListboxButton>
    <ListboxOptions class="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-56 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none text-sm">
      <ListboxOption as="template"
                     v-for="record in recordTypes"
                     :key="record"
                     :value="record"
                     v-slot="{ active, selected }">
        <li :class="[active ? 'text-white bg-indigo-600' : 'text-gray-900', 'cursor-default select-none relative py-2 pl-3 pr-9']">
          <div class="flex items-center">
            <span :class="[selected ? 'font-semibold' : 'font-normal', 'ml-3 block truncate']">
              {{ record }}
            </span>
          </div>
        </li>
      </ListboxOption>
    </ListboxOptions>
  </Listbox>
</template>

<script>
import {
  Listbox,
  ListboxButton,
  ListboxOptions,
  ListboxOption,
} from "@headlessui/vue";
import { SelectorIcon } from "@heroicons/vue/solid";

const recordTypes = ["A", "AAAA", "CNAME"];

export default {
  components: {
    Listbox,
    ListboxButton,
    ListboxOptions,
    ListboxOption,
    SelectorIcon,
  },
  props: {
    modelValue: String,
  },
  emits: ["update:modelValue"],
  data () {
    return { recordTypes }
  }
};
</script>
<template>
  <div class="header wrapper">
    <UiToastList/>
  </div>

</template>

<script>
import UiToastList from '@/components/UiToastList'

export default {
  name: 'Header',
  data: () => ({
    title: 'Route Simulator'
  }),
  components: {
    UiToastList
  },
  methods: {}
}
</script>

<style lang="scss" scoped>
.header {
  z-index: 1;
  padding: 20px 0;
  position: relative;

  .is-active {
    color: #e01b3c;
  }

  a {
    color: #000;
    text-decoration: none;
  }
}
</style>

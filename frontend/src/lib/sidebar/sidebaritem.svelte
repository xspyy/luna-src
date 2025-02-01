<script>
    export let extension = ".lua";
    export let classDecision = "";
</script>

<div
    style="width: calc(100% - 10px); height: 50px; display: flex; gap: 7px; padding-left: 10px; align-items: center; color: #FFFFFFaa; border-bottom: solid 1.5px #2a2831;"
    class={classDecision + " tab-item"}
>
    <img class={classDecision} alt="search" src="file.svg" />
    <slot/>
    <label class={classDecision} style="opacity: 50%; font-weight: 400; margin-left: -6px;-webkit-touch-callout: none;
  -webkit-user-select: none;
  -khtml-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;"
        >
        {extension}
    </label>
</div>

<style>
    .active-container {
        background: linear-gradient(
            0deg,
            rgba(85, 85, 85, 0.09287464985994398) 0%,
            rgba(94, 94, 94, 0.06766456582633051) 47%,
            rgba(99, 99, 99, 0.09007352941176472) 100%
        );
    }
    .tab-item:hover {
        background: linear-gradient(
            0deg,
            rgba(85, 85, 85, 0.09287464985994398) 0%,
            rgba(94, 94, 94, 0.06766456582633051) 47%,
            rgba(99, 99, 99, 0.09007352941176472) 100%
        );
    }
</style>

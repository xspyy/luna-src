<script>
    export let iconpath = "editor.svg";
    export let classCustom = "";
    let iconname = iconpath.replace(".svg", "");
    export let MarginTop = true;
    export let customContent = false;
    $: classMarginTop = MarginTop ? "margin-top: 10px;" : "";
    export let DontTouch = true;

    $: classDontTouch = DontTouch ? "contain:content;-webkit-touch-callout: none;-webkit-user-select: none;-khtml-user-select: none;-moz-user-select: none;-ms-user-select: none;" : "";

    export let onclicked;
    
</script>

<button style={classMarginTop + classDontTouch} onclick={onclicked} class={classCustom}>
    {#if customContent} 
    <slot/>
    {:else}
    <img alt={iconname} src={iconpath}/>
    {/if}
    
</button>

<style>
    button {
        align-items: center;
        border-radius: 10px;
        justify-content: center;
        display: flex;
        padding: 10px;
        background: radial-gradient(100% 100% at 51.43% 100%, #2a2931c5 20%, rgba(42, 41, 49, 0.151) 100%);
        border-width: 0px;
        transition: 150ms;
        color: #ffffff9e; 
        font-family: 'Outfit';
        border: solid 1px rgba(108, 103, 128, 0.096);
        transition: transform 0.2s ease-in-out;
    }
    

    button:focus { outline: none; }

    button:hover {
        background: radial-gradient(100% 100% at 51.43% 100%, #2A2931 20%, rgba(42, 41, 49, 0.00) 100%);
        backdrop-filter: blur(15px);
        border: solid 1px rgba(108, 103, 128, 0.25);
        position: relative;transform: scale(1.014);
    }

    .active {
        
        background: radial-gradient(100% 100% at 51.43% 100%, #2A2931 20%, rgba(42, 41, 49, 0.00) 100%);
        backdrop-filter: blur(15px);
        border: solid 1px rgba(108, 103, 128, 0.25);
        position: relative;
        
    }
</style>

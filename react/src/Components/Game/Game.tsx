import * as React from 'react'
import { styled } from '@mui/material/styles'
import Card from '@mui/material/Card'
import CardMedia from '@mui/material/CardMedia'
import CardContent from '@mui/material/CardContent'
import CardActions from '@mui/material/CardActions'
import IconButton, { IconButtonProps } from '@mui/material/IconButton'
import Typography from '@mui/material/Typography'
import TuneIcon from '@mui/icons-material/Tune'
import EditIcon from '@mui/icons-material/Edit'
import MapsHomeWorkIcon from '@mui/icons-material/MapsHomeWork'
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline'
import PlayCircleOutlineIcon from '@mui/icons-material/PlayCircleOutline'
import { AppBar, ButtonGroup, ToggleButton, ToggleButtonGroup } from '@mui/material'
import { Badge } from '@mui/material'
import { BadgeProps } from '@mui/material'
import { Box } from '@mui/material'
import { Button } from '@mui/material'
import { Dialog } from '@mui/material'
import { DialogActions } from '@mui/material'
import { DialogContent } from '@mui/material'
import { DialogTitle } from '@mui/material'
import { FormControl } from '@mui/material'
import { Grid } from '@mui/material'
import { InputAdornment } from '@mui/material'
import { InputLabel } from '@mui/material'
import { ListItemIcon } from '@mui/material'
import { Menu } from '@mui/material'
import { MenuItem } from '@mui/material'
import { Rating } from '@mui/material'
import { Select } from '@mui/material'
import { TextField } from '@mui/material'
import { Toolbar } from '@mui/material'
import { Tooltip } from '@mui/material'
import { Divider } from '@mui/material'
import { FormControlLabel } from '@mui/material'
import { FormGroup } from '@mui/material'
import { InputBase } from '@mui/material'
import { List } from '@mui/material'
import { ListItem } from '@mui/material'
import { ListItemText } from '@mui/material'
import { Paper } from '@mui/material'
import { Switch } from '@mui/material'
import { Tabs } from '@mui/material'
import { alpha } from '@mui/material'
import Pagination from '@mui/material/Pagination'
import Tab from '@mui/material/Tab'
import TabContext from '@mui/lab/TabContext'
import TabPanel from '@mui/lab/TabPanel'
import LinearProgress, { LinearProgressProps } from '@mui/material/LinearProgress'
import { Code, CodeSlash, CodeSquare, Heart, HeartFill, HeartPulse, PatchCheck, PatchExclamation, PatchQuestion } from 'react-bootstrap-icons'
import { Fragment, useEffect, useState } from 'react'
import PostAddIcon from '@mui/icons-material/PostAdd'
import FormatListNumberedIcon from '@mui/icons-material/FormatListNumbered'
import Snackbar from '@mui/material/Snackbar'
import MuiAlert, { AlertProps } from '@mui/material/Alert'
import SearchIcon from '@mui/icons-material/Search'
import { purple } from '@mui/material/colors'
import MenuIcon from '@mui/icons-material/Menu'
import AlarmOffIcon from '@mui/icons-material/AlarmOff'
import useMediaQuery from '@mui/material/useMediaQuery'
import { GameData, StopwatchData } from '../../interfaces'
import selectOptions from '../../selectOptions.json'
import Cookies from 'js-cookie'

// Icons
import { Tablet } from 'react-bootstrap-icons'
import { PcDisplay } from 'react-bootstrap-icons'
import { NintendoSwitch } from 'react-bootstrap-icons'
import { Playstation } from 'react-bootstrap-icons'
import { Xbox } from 'react-bootstrap-icons'
import { CheckSquare } from 'react-bootstrap-icons'
import { Square } from 'react-bootstrap-icons'
import { getCSRFToken, hourOfDuration, minOfDuration, percentage, rating } from '../../utils'
import { DoingBadge, DoneBadge, TodoBadge } from '../Common/Badges'
import TodoIcon from '@mui/icons-material/SyncLock'
import DoingIcon from '@mui/icons-material/Sync'
import DoneIcon from '@mui/icons-material/PublishedWithChanges'

export default function Game(props: { authed: boolean }) {
    // const csrftoken = Cookies.get('csrftoken')
    // console.log(csrftoken)

    const api = process.env.API_URL;

    const [csrfToken, setCsrfToken] = useState<string>('');

    const [stopwatch, setStopwatch] = useState<StopwatchData>()

    const pcScreen = useMediaQuery('(min-width:600px)')
    function matchScreen(pc: number, mobile: number): number {
        if (pcScreen === true) {
            return pc
        } else {
            return mobile
        }
    }

    // States
    const [games, setGames] = useState<GameData[]>([])

    const [tabPlatform, setTabPlatform] = useState('All')
    const platformColor = (platform: string) => {
        switch (platform) {
            case 'All':
                return '#FFF176'
            case 'PC':
                return '#E3963E'
            case 'PlayStation':
                return '#2E6DB4'
            case 'Nintendo Switch':
                return '#E60012'
            case 'Xbox':
                return '#107C10'
            case 'Mobile':
                return '#730073'
            default:
                return 'gray'
        }
    }

    const [page, setPage] = useState(1)
    const [totalPages, setTotalPages] = useState(1)
    const [tabStatus, setTabStatus] = useState("Playing")

    // Status counts
    const [playedCnt, setPlayedCnt] = useState(0)
    const [playingCnt, setPlayingCnt] = useState(0)
    const [toPlayCnt, setToPlayCnt] = useState(0)

    // Platform counts
    const [allCount, setAllCount] = useState(0)
    const [pcCount, setPcCount] = useState(0)
    const [psCount, setPsCount] = useState(0)
    const [nsCount, setNsCount] = useState(0)
    const [xboxCount, setXboxCount] = useState(0)
    const [mobileCount, setMobileCount] = useState(0)

    // Effects
    useEffect(() => {
        const getCsrfToken = async () => {
            try {
                const response = await fetch('/api/csrf/');
                const data = await response.json();
                setCsrfToken(data['csrf_token']);
            } catch (error) {
                console.error('Error fetching CSRF token:', error);
            }
        };

        getCsrfToken();

        refresh()
    }, [tabStatus, tabPlatform, page])

    const SearchBar = styled('div')(({ theme }) => ({
        position: 'relative',
        borderRadius: theme.shape.borderRadius,
        backgroundColor: alpha(theme.palette.common.white, 0.15),
        '&:hover': {
            backgroundColor: alpha(theme.palette.common.white, 0.25),
        },
        marginLeft: 0,
        width: '100%',
        [theme.breakpoints.up('sm')]: {
            marginLeft: theme.spacing(1),
            width: 'auto',
        },
    }))

    const SearchIconWrapper = styled('div')(({ theme }) => ({
        padding: theme.spacing(0, 2),
        height: '100%',
        position: 'absolute',
        pointerEvents: 'none',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
    }))

    const StyledInputBase = styled(InputBase)(({ theme }) => ({
        color: 'inherit',
        '& .MuiInputBase-input': {
            padding: theme.spacing(1, 0.7, 1, 0),
            // vertical padding + font size from searchIcon
            paddingLeft: `calc(1em + ${theme.spacing(4)})`,
            transition: theme.transitions.create('width'),
            width: '100%',
            [theme.breakpoints.up('sm')]: {
                width: '0ch',
                '&:focus': {
                    width: '10ch',
                },
            },
        },
    }))

    function Progressing(props: LinearProgressProps & { value: number }) {
        return (
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Box sx={{ width: '100%' }}>
                    {LinearProgressing(props)}
                </Box>
                <Box sx={{ minWidth: 35, pl: 1 }}>
                    {ProgressingPercentage(props.value)}
                </Box>
            </Box>
        )
    }

    function LinearProgressing(props: LinearProgressProps & { value: number }) {
        switch (props.value) {
            case -1:
                return <LinearProgress variant="determinate" color="info" {...props} />
            case 0:
                return <LinearProgress variant="determinate" color="secondary" {...props} />
            case 100:
                return <LinearProgress variant="determinate" color="inherit" {...props} />
            default:
                return <LinearProgress variant="determinate" color="success" {...props} />
        }
    }

    function ProgressingPercentage(value: number) {
        switch (value) {
            case -1:
                return <Typography variant="body2" sx={{ color: "info.main" }}>♾️</Typography>
            case 0:
                return <Typography variant="body2" sx={{ color: "secondary.main" }}>{value}%</Typography>
            case 100:
                return <Typography variant="body2" sx={{ color: "inhert.main" }}>{value}%</Typography>
            default:
                return <Typography variant="body2" sx={{ color: "success.main" }}>{value}%</Typography>
        }
    }

    // Create game dialog
    const [openCreateGameDialog, setOpenCreateGameDialog] = useState(false);
    const handleCreateGameDialogOpen = () => { setOpenCreateGameDialog(true); }
    const handleCreateGameDialogClose = () => {
        setOpenCreateGameDialog(false)
    }

    // Update Game Dialog
    // const handleUpdateGameSubmit = async (event: React.FormEvent) => {
    //     event.preventDefault();
    //     const formData = new FormData(event.target as HTMLFormElement);

    //     const response = await fetch(`/api/game/update/${updateGame?.id}/`, {
    //         method: 'POST',
    //         headers: {
    //             'X-CSRFToken': csrftoken || '',
    //         },
    //         body: formData,
    //     });

    //     if (response.ok) {
    //         const data = await response.json();
    //         console.log('Success:', data);
    //     } else {
    //         console.error('Failed to submit form');
    //     }
    // }
    const [openUpdateGameDialog, setOpenUpdateGameDialog] = useState(false)
    const handleUpdateGameDialogOpen = (id: String) => { fetchUpdateGame(id) }
    const handleUpdateGameDialogClose = () => {
        setOpenUpdateGameDialog(false)
    }

    const [successAlertOpen, setSuccessAlertOpen] = useState(false)
    const [warningAlertOpen, setWarningAlertOpen] = useState(false)
    const [errorAlertOpen, setErrorAlertOpen] = useState(false)

    const Alert = React.forwardRef<HTMLDivElement, AlertProps>(
        function Alert(props, ref) {
            return <MuiAlert elevation={6} ref={ref} variant="filled" {...props} />
        }
    )

    const handleSuccessAlertOpen = () => { setSuccessAlertOpen(true) }
    const handleSuccessAlertClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') { return }
        setSuccessAlertOpen(false)
    }

    const handleWarningAlertOpen = () => { setWarningAlertOpen(true) }
    const handleWarningAlertClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') { return }
        setWarningAlertOpen(false)
    }

    const handleErrorAlertOpen = () => { setErrorAlertOpen(true) }
    const handleErrorAlertClose = (event?: React.SyntheticEvent | Event, reason?: string) => {
        if (reason === 'clickaway') { return }
        setErrorAlertOpen(false)
    }

    const [updateGame, setUpdateGame] = useState<GameData>()

    function fetchUpdateGame(id: String) {
        fetch(`/api/game/update?id=${id}`, { method: "GET" })
            .then(resp => resp.json())
            .then(data => {
                if (data != null) {
                    setUpdateGame(data["game"])
                    setOpenUpdateGameDialog(true)
                }
            });
    }

    function fetchStopwatch() {
        fetch(`/api/game/stopwatch`, { method: "GET" })
            .then(resp => resp.json())
            .then(data => {
                if (data != null) {
                    setStopwatch(data["stopwatch"])
                }
            })
            .catch(error => console.log(error))
    }

    function refresh() {
        fetchStopwatch()

        fetch(`/api/game/pages?status=${tabStatus}&platform=${tabPlatform}&p=${page}`)
            .then(resp => resp.json())
            .then(data => {
                if (data["games"] != null) {
                    setGames(data["games"])
                } else {
                    setGames([])
                }
                setTotalPages(data["pages"]);
            })

        fetch(`/api/game/badge?status=${tabStatus}`)
            .then(resp => resp.json())
            .then(data => {
                setPlayedCnt(data["badges"]["played"]);
                setPlayingCnt(data["badges"]["playing"]);
                setToPlayCnt(data["badges"]["to_play"]);

                setAllCount(data["badges"]["all_platform"]);
                setPcCount(data["badges"]["pc"]);
                setPsCount(data["badges"]["playstation"]);
                setNsCount(data["badges"]["nintendo_switch"]);
                setXboxCount(data["badges"]["xbox"]);
                setMobileCount(data["badges"]["mobile"]);
            })
    }

    const handleDeleteGame = (id: string | undefined) => {
        fetch(`/api/game/delete?id=${id}`, {
            method: "DELETE",
        })
            .then(() => {
                handleUpdateGameDialogClose();
                refresh()
            })
    }

    const handlePageChange = (event: React.ChangeEvent<unknown>, value: number) => {
        setPage(value)
    }

    const handleStatusChange = (
        event: React.MouseEvent<HTMLElement>,
        newValue: string,
    ) => {
        if (newValue === null) {
            return
        } else {
            setPage(1)
            setTabStatus(newValue)
        }
    }

    const handlePlatformChange = (event: React.SyntheticEvent, newValue: string) => {
        setPage(1)
        setTabPlatform(newValue)
    };

    const handleStartGaming = (id: string) => {
        fetch(`/api/game/start?id=${id}`, { method: "GET" })
            .then(resp => {
                if (!resp.ok) {
                    handleErrorAlertOpen()
                } else {
                    refresh()
                    handleSuccessAlertOpen()
                }
            })
    }

    const handleStopStopwatch = () => {
        fetch(`/api/game/stop/`, { method: "GET" })
            .then(resp => {
                handleWarningAlertOpen()
                refresh()
            })
    }

    const StyledRating = styled(Rating)({
        '& .MuiRating-iconFilled': {
            color: '#ff6d75',
        },
        '& .MuiRating-iconHover': {
            color: '#ff3d47',
        },
    })

    const StyledBadge = styled(Badge)<BadgeProps>(({ theme }) => ({
        '& .MuiBadge-badge': {
            right: -3,
            top: 2,
            border: `2px solid ${theme.palette.background.paper}`,
            padding: '0 2px',
        },
    }))

    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null)
    const open = Boolean(anchorEl)
    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget)
    }
    const handleClose = () => { setAnchorEl(null) }

    return (
        <Fragment>
            <Box
                sx={{
                    m: matchScreen(4, 1),
                    p: 1,
                    border: 2,
                    borderRadius: 1,
                    borderColor: 'divider'
                }}
            >
                <AppBar position='static'>
                    <Toolbar>
                        <Box sx={{ display: 'flex', alignItems: 'center', textAlign: 'center' }}>
                            <IconButton
                                onClick={handleClick}
                                aria-controls={open ? 'account-menu' : undefined}
                                aria-haspopup="true"
                                aria-expanded={open ? 'true' : undefined}
                            >
                                <MenuIcon
                                    sx={{
                                        fontSize: 32,
                                        color: purple[100],
                                        "&:hover": { color: purple[200], fontSize: 35 }
                                    }}
                                />
                            </IconButton>
                        </Box>
                        <Menu
                            anchorEl={anchorEl}
                            id="account-menu"
                            open={open}
                            onClose={handleClose}
                            onClick={handleClose}
                            transformOrigin={{ horizontal: 'right', vertical: 'top' }}
                            anchorOrigin={{ horizontal: 'right', vertical: 'bottom' }}
                        >
                            {props.authed ?
                                <MenuItem onClick={handleCreateGameDialogOpen}>
                                    <ListItemIcon>
                                        <PostAddIcon
                                            sx={{
                                                fontSize: 28,
                                                color: purple[100],
                                                "&:hover": { color: purple[200] }
                                            }}
                                        />
                                    </ListItemIcon>
                                    <Typography
                                        sx={{
                                            pl: 1,
                                            color: purple[100]
                                        }}
                                    >
                                        Create Game
                                    </Typography>
                                </MenuItem>
                            :<></>}

                            <MenuItem onClick={handleClose}>
                                <ListItemIcon>
                                    <FormatListNumberedIcon
                                        sx={{
                                            fontSize: 28,
                                            color: purple[100],
                                            "&:hover": { color: purple[200] }
                                        }}
                                    />
                                </ListItemIcon>
                                <Typography
                                    sx={{
                                        pl: 1,
                                        color: purple[100]
                                    }}
                                >
                                    Ranking
                                </Typography>
                            </MenuItem>
                        </Menu>

                        <ToggleButtonGroup
                            exclusive
                            sx={{ ml: 1, mt: 1, mb: 1 }}
                            value={tabStatus}
                            onChange={handleStatusChange}
                        >
                            <ToggleButton value="Played">
                                {DoneBadge(playedCnt)}
                            </ToggleButton>
                            <ToggleButton value="Playing">
                                {DoingBadge(playingCnt)}
                            </ToggleButton>
                            <ToggleButton value="ToPlay">
                                {TodoBadge(toPlayCnt)}
                            </ToggleButton>
                        </ToggleButtonGroup>

                        <Box display='flex' flexGrow={1} />

                        <SearchBar>
                            <SearchIconWrapper>
                                <SearchIcon />
                            </SearchIconWrapper>
                            <StyledInputBase />
                        </SearchBar> 
                    </Toolbar>
                </AppBar>

                <TabContext value={tabStatus}>
                    <TabPanel value={tabStatus}>
                        <Grid container spacing={2}>
                            <Grid
                                item
                                style={{ width: '100px' }}
                                sx={{
                                    ml: -5,
                                    borderRight: 1,
                                    borderColor: 'divider'
                                }}
                            >
                                <Tabs
                                    variant="fullWidth"
                                    orientation="vertical"
                                    TabIndicatorProps={{
                                        style: { background: platformColor(tabPlatform) }
                                    }}
                                    value={tabPlatform}
                                    onChange={handlePlatformChange}
                                >
                                    <Tab
                                        value="All"
                                        icon={
                                            <StyledBadge
                                                badgeContent={allCount}
                                                sx={{
                                                    "& .MuiBadge-badge": {
                                                        color: 'black',
                                                        backgroundColor: platformColor(tabPlatform === 'All' ? tabPlatform : '')
                                                    }
                                                }}
                                            >
                                                {tabPlatform === 'All'
                                                    ? <CheckSquare color='#FFF176' size={35} />
                                                    : <Square color='gray' size={35} />}
                                            </StyledBadge>
                                        }
                                    />

                                    <Tab
                                        value="PC"
                                        sx={{ mt: 1 }}
                                        icon={
                                            <StyledBadge
                                                badgeContent={pcCount}
                                                sx={{
                                                    "& .MuiBadge-badge": {
                                                        color: 'black',
                                                        backgroundColor: platformColor(tabPlatform === 'PC' ? tabPlatform : '')
                                                    }
                                                }}
                                            >
                                                <PcDisplay color={platformColor(tabPlatform === 'PC' ? tabPlatform : '')} size={35} />
                                            </StyledBadge>
                                        }
                                    />

                                    <Tab
                                        value="PlayStation"
                                        sx={{ mt: 1 }}
                                        icon={
                                            <StyledBadge
                                                badgeContent={psCount}
                                                sx={{
                                                    "& .MuiBadge-badge": {
                                                        color: 'black',
                                                        backgroundColor: platformColor(tabPlatform === 'PlayStation' ? tabPlatform : '')
                                                    }
                                                }}
                                            >
                                                <Playstation color={platformColor(tabPlatform === 'PlayStation' ? tabPlatform : '')} size={35} />
                                            </StyledBadge>
                                        }
                                    />

                                    <Tab
                                        value="Nintendo Switch"
                                        sx={{ mt: 1 }}
                                        icon={
                                            <StyledBadge
                                                badgeContent={nsCount}
                                                sx={{
                                                    "& .MuiBadge-badge": {
                                                        color: 'black',
                                                        backgroundColor: platformColor(tabPlatform === 'Nintendo Switch' ? tabPlatform : '')
                                                    }
                                                }}
                                            >
                                                <NintendoSwitch color={platformColor(tabPlatform === 'Nintendo Switch' ? tabPlatform : '')} size={35} />
                                            </StyledBadge>
                                        }
                                    />

                                    <Tab
                                        value="Xbox"
                                        sx={{ mt: 1 }}
                                        icon={
                                            <StyledBadge
                                                badgeContent={xboxCount}
                                                sx={{
                                                    "& .MuiBadge-badge": {
                                                        color: 'black',
                                                        backgroundColor: platformColor(tabPlatform === 'Xbox' ? tabPlatform : '')
                                                    }
                                                }}
                                            >
                                                <Xbox color={platformColor(tabPlatform === 'Xbox' ? tabPlatform : '')} size={35} />
                                            </StyledBadge>
                                        }
                                    />

                                    <Tab
                                        value="Mobile"
                                        sx={{ mt: 1 }}
                                        icon={
                                            <StyledBadge
                                                sx={{
                                                    "& .MuiBadge-badge": {
                                                        color: 'black',
                                                        backgroundColor: platformColor(tabPlatform === 'Mobile' ? tabPlatform : '')
                                                    }
                                                }}
                                                badgeContent={mobileCount}
                                            >
                                                <Tablet color={platformColor(tabPlatform === 'Mobile' ? tabPlatform : '')} size={35} />
                                            </StyledBadge>
                                        }
                                    />
                                </Tabs>
                            </Grid>

                            <Grid
                                item
                                xs={true}
                                justifyContent='center'
                                alignItems='center'
                            >
                                <Grid container spacing={3}>
                                    {games?.map((g, i) => (
                                        <Grid item>
                                            <Card
                                                sx={{
                                                    maxWidth: 250,
                                                    maxHeight: 565,
                                                }}
                                                key={g.id}
                                            >
                                                <CardMedia
                                                    component="img"
                                                    height="250"
                                                    image={"/assets/images/games/" + g.id + ".webp"}
                                                />
                                                <CardContent>
                                                    <Tooltip title="">
                                                        <Typography
                                                            variant="subtitle1"
                                                            align="center"
                                                            color="text.info"
                                                            sx={{ height: 30 }}
                                                        >
                                                            {g.title}
                                                        </Typography>
                                                    </Tooltip>
                                                </CardContent>

                                                {props.authed ?
                                                    <CardActions disableSpacing>
                                                        <Tooltip title="Property">
                                                            <IconButton onClick={() => handleUpdateGameDialogOpen(g.id)}>
                                                                <TuneIcon />
                                                            </IconButton>
                                                        </Tooltip>

                                                        {stopwatch?.game_id == g.id ?
                                                            <Tooltip sx={{ ml: -1.5 }} title="Stop">
                                                                <IconButton onClick={() => handleStopStopwatch()}>
                                                                    <AlarmOffIcon color="warning" />
                                                                </IconButton>
                                                            </Tooltip>
                                                            :
                                                            <Tooltip sx={{ ml: -1.5 }} title="Start">
                                                                <IconButton onClick={() => handleStartGaming(g.id)}>
                                                                    <PlayCircleOutlineIcon color="success" />
                                                                </IconButton>
                                                            </Tooltip>}
                                                    </CardActions>
                                                :<></>}

                                                <CardContent sx={{ mt: -4 }}>
                                                    <Box
                                                        sx={{
                                                            mx: "auto",
                                                            display: 'flex',
                                                            flexDirection: 'column',
                                                            alignItems: 'center',
                                                            '& > *': { m: 1 },
                                                        }}
                                                    >
                                                        <TextField
                                                            fullWidth
                                                            size="small"
                                                            sx={{ pt: 0.5 }}
                                                            inputProps={{
                                                                style: { textAlign: 'right' },
                                                                readOnly: true,
                                                            }}
                                                            value={rating(g.rating)}
                                                            InputProps={{
                                                                startAdornment: (
                                                                    <InputAdornment position="start">
                                                                        {g.platform === 'Mobile' ? <Tablet /> : <></>}
                                                                        {g.platform === 'Nintendo Switch' ? <NintendoSwitch /> : <></>}
                                                                        {g.platform === 'PC' ? <PcDisplay /> : <></>}
                                                                        {g.platform === 'PlayStation' ? <Playstation /> : <></>}
                                                                        {g.platform === 'Xbox' ? <Xbox /> : <></>}
                                                                    </InputAdornment>
                                                                )
                                                            }}
                                                        />

                                                        <TextField
                                                            fullWidth
                                                            size="small"
                                                            sx={{ pt: 0.5 }}
                                                            inputProps={{
                                                                style: { textAlign: 'right' },
                                                                readOnly: true,
                                                            }}
                                                            value={g.developer}
                                                            InputProps={{
                                                                startAdornment: (
                                                                    <InputAdornment position="start">
                                                                        <Code />
                                                                    </InputAdornment>
                                                                )
                                                            }}
                                                        />

                                                        <TextField
                                                            fullWidth
                                                            size="small"
                                                            sx={{ pt: 0.5 }}
                                                            inputProps={{ style: { textAlign: 'right' }, readOnly: true }}
                                                            value={g.publisher}
                                                            InputProps={{
                                                                startAdornment: (
                                                                    <InputAdornment position="start"><CodeSquare /></InputAdornment>
                                                                )
                                                            }}
                                                        />

                                                        <TextField
                                                            fullWidth
                                                            size="small"
                                                            sx={{ pt: 0.5 }}
                                                            inputProps={{ style: { textAlign: 'right' }, readOnly: true }}
                                                            value={hourOfDuration(g.played_time)}
                                                            InputProps={{
                                                                startAdornment: (
                                                                    <InputAdornment position="start">
                                                                        {g.status === 'Played' ? <HeartFill /> : <></>}
                                                                        {g.status === 'Playing' ? <HeartPulse /> : <></>}
                                                                        {g.status === 'ToPlay' ? <Heart /> : <></>}
                                                                    </InputAdornment>
                                                                ),
                                                                endAdornment: (
                                                                    <InputAdornment position="end">Hour(s)</InputAdornment>
                                                                )
                                                            }}
                                                        />

                                                        <Box sx={{ width: '100%' }}>
                                                            <Progressing value={percentage(g.played_time, g.time_to_beat)} />
                                                        </Box>
                                                    </Box>
                                                </CardContent>
                                            </Card>
                                        </Grid>
                                    ))}
                                </Grid>

                                <Box
                                    sx={{ width: '100%', pt: 3, pb: 3 }}
                                    display="flex"
                                    justifyContent="center"
                                    alignItems="center"
                                >
                                    {totalPages > 1 ?
                                        <Pagination
                                            count={totalPages}
                                            page={page}
                                            onChange={handlePageChange}
                                            variant="outlined"
                                            color="secondary"
                                        />
                                        : <></>}
                                </Box>
                            </Grid>
                        </Grid>


                    </TabPanel>

                </TabContext>

                <Snackbar open={successAlertOpen} autoHideDuration={3000} onClose={handleSuccessAlertClose}>
                    <Alert severity="success" onClose={handleSuccessAlertClose} sx={{ width: '100%' }}>Started!</Alert>
                </Snackbar>
                <Snackbar open={warningAlertOpen} autoHideDuration={3000} onClose={handleWarningAlertClose}>
                    <Alert severity="warning" onClose={handleWarningAlertClose} sx={{ width: '100%' }}>Closed!</Alert>
                </Snackbar>
                <Snackbar open={errorAlertOpen} autoHideDuration={3000} onClose={handleErrorAlertClose}>
                    <Alert severity="error" onClose={handleErrorAlertClose} sx={{ width: '100%' }}>Failed!</Alert>
                </Snackbar>

                <Dialog
                    open={openUpdateGameDialog}
                    onClose={handleUpdateGameDialogClose}
                >
                    <DialogTitle align="center">
                        Update Game
                    </DialogTitle>
                    <DialogContent>
                        <form method="post" encType="multipart/form-data" action="/api/game/update/">
                        {/* <form method="post" encType="multipart/form-data" action={`/api/game/update/${updateGame?.id}/`} onSubmit={handleUpdateGameSubmit}> */}
                            <input type="hidden" name="csrfmiddlewaretoken" value={csrfToken || ''} />

                            <FormControl fullWidth sx={{ mt: 1 }}>
                                <TextField
                                    name="id"
                                    label="Id"
                                    defaultValue={updateGame?.id}
                                    inputProps={{
                                        readOnly: true
                                    }}
                                />
                            </FormControl>

                            <FormControl fullWidth sx={{ mt: 2 }}>
                                <TextField
                                    name="title"
                                    label="Title"
                                    defaultValue={updateGame?.title}
                                >
                                </TextField>
                            </FormControl>

                            <Grid container>
                                <Grid item sx={{ width: '48%' }}>
                                    <FormControl fullWidth sx={{ mt: 2 }}>
                                        <TextField
                                            name="developer"
                                            label="Developer"
                                            defaultValue={updateGame?.developer}
                                        />
                                    </FormControl>
                                </Grid>

                                <Grid item sx={{ width: '4%' }} />

                                <Grid item sx={{ width: '48%' }}>
                                    <FormControl fullWidth sx={{ mt: 2 }}>
                                        <TextField
                                            name="publisher"
                                            label="Publisher"
                                            defaultValue={updateGame?.publisher}
                                        />
                                    </FormControl>
                                </Grid>
                            </Grid>

                            <Grid container>
                                <Grid item sx={{ width: '48%' }}>
                                    <FormControl
                                        fullWidth
                                        sx={{ mt: 2 }}
                                    >
                                        <InputLabel htmlFor="Status">Status</InputLabel>
                                        <Select
                                            name="status"
                                            label="Status"
                                            defaultValue={updateGame?.status}
                                        >
                                            <MenuItem key="Played" value="Played">Played</MenuItem>
                                            <MenuItem key="Playing" value="Playing">Playing</MenuItem>
                                            <MenuItem key="ToPlay" value="ToPlay">ToPlay</MenuItem>
                                        </Select>
                                    </FormControl>
                                </Grid>

                                <Grid item sx={{ width: '4%' }} />

                                <Grid item sx={{ width: '48%' }}>
                                    <FormControl
                                        fullWidth
                                        sx={{ mt: 2 }}
                                    >
                                        <InputLabel htmlFor="Genre">Genre</InputLabel>
                                        <Select
                                            name="genre"
                                            label="Genre"
                                            defaultValue={updateGame?.genre}
                                        >
                                            {selectOptions.genres.map((option: any, index) => {
                                                return (
                                                    <MenuItem key={index} value={option.value}>{option.label}</MenuItem>
                                                )
                                            })}
                                        </Select>
                                    </FormControl>
                                </Grid>
                            </Grid>

                            <Grid container>
                                <Grid item sx={{ width: '48%' }}>
                                    <FormControl
                                        fullWidth
                                        sx={{ mt: 2 }}
                                    >
                                        <InputLabel htmlFor="Platform">Platform</InputLabel>
                                        <Select
                                            name="platform"
                                            label="Platform"
                                            defaultValue={updateGame?.platform}
                                        >
                                            {selectOptions.platforms.map((option: any, index) => {
                                                return (
                                                    <MenuItem key={index} value={option.value}>{option.label}</MenuItem>
                                                )
                                            })}
                                        </Select>
                                    </FormControl>
                                </Grid>
                                <Grid item sx={{ width: '4%' }}></Grid>
                                <Grid item sx={{ width: '48%' }}>
                                    <FormControl fullWidth sx={{ mt: 2 }}>
                                        <InputLabel htmlFor="Rating">Rating</InputLabel>
                                        <Select
                                            name="rating"
                                            label="Rating"
                                            defaultValue={updateGame?.rating}
                                        >
                                            {selectOptions.rating.map((option: any, index) => {
                                                return (
                                                    <MenuItem key={index} value={option.value}>{option.label}</MenuItem>
                                                )
                                            })}
                                        </Select>
                                    </FormControl>
                                </Grid>
                            </Grid>

                            <Grid container sx={{ mt: 2 }}>
                                <Grid item sx={{ width: '30%' }}>
                                    <FormControl fullWidth>
                                        <TextField
                                            name="played_time_hour"
                                            type="number"
                                            label="Played Hour"
                                            defaultValue={hourOfDuration(updateGame?.played_time)}
                                        >
                                        </TextField>
                                    </FormControl>
                                </Grid>
                                <Grid item sx={{ width: '1%' }}></Grid>
                                <Grid item sx={{ width: '17%' }}>
                                    <FormControl fullWidth>
                                        <TextField
                                            name="played_time_min"
                                            type="number"
                                            label="Min"
                                            defaultValue={minOfDuration(updateGame?.played_time)}
                                        >
                                        </TextField>
                                    </FormControl>
                                </Grid>

                                <Grid item sx={{ width: '4%' }}></Grid>

                                <Grid item sx={{ width: '30%' }}>
                                    <FormControl fullWidth>
                                        <TextField
                                            name="time_to_beat_hour"
                                            type="number"
                                            label="Hour to Beat"
                                            defaultValue={hourOfDuration(updateGame?.time_to_beat)}
                                        >
                                        </TextField>
                                    </FormControl>
                                </Grid>

                                <Grid item sx={{ width: '1%' }}></Grid>

                                <Grid item sx={{ width: '17%' }}>
                                    <FormControl fullWidth>
                                        <TextField
                                            name="time_to_beat_min"
                                            type="number"
                                            label="Min"
                                            defaultValue={minOfDuration(updateGame?.time_to_beat)}
                                        >
                                        </TextField>
                                    </FormControl>
                                </Grid>
                            </Grid>

                            <Grid container sx={{ mt: 1 }} justifyContent="flex-end">
                                <Grid item xs={12}>
                                    <FormControl >
                                        <input type="file" id="cover" name="cover" />
                                    </FormControl>
                                </Grid>
                            </Grid>

                            <DialogActions style={{ justifyContent: "space-between" }} sx={{ mt: 1, mb: -1, ml: -1, mr: -1 }}>
                                <Button color="error" onClick={e => handleDeleteGame(updateGame?.id)}>Delete</Button>
                                <Box>
                                    <Button color="secondary" onClick={handleUpdateGameDialogClose}>Cancel</Button>
                                    <Button color="success" type="submit">Update</Button>
                                </Box>
                            </DialogActions>
                        </form>
                    </DialogContent>
                </Dialog>
            </Box>

            <Dialog
                open={openCreateGameDialog}
                onClose={handleCreateGameDialogClose}
            >
                <DialogTitle align="center">Create Game</DialogTitle>
                <DialogContent>
                    <form method="post" action="/api/game/create">
                        <FormControl fullWidth sx={{ mt: 1 }}>
                            <TextField name="title" label="Title" required />
                        </FormControl>

                        <FormControl fullWidth sx={{ mt: 2 }}>
                            <TextField name="developer" label="Developer" required />
                        </FormControl>

                        <FormControl fullWidth sx={{ mt: 2 }}>
                            <TextField name="publisher" label="Publisher" required />
                        </FormControl>

                        <FormControl fullWidth sx={{ mt: 2 }}>
                            <InputLabel htmlFor="Genre">Genre</InputLabel>
                            <Select name="genre" label="Genre" required>
                                {selectOptions.genres.map((option: any, index) => {
                                    return (
                                        <MenuItem key={index} value={option.value}>{option.label}</MenuItem>
                                    )
                                })}
                            </Select>
                        </FormControl>

                        <FormControl fullWidth sx={{ mt: 2 }}>
                            <InputLabel htmlFor="Platform">Platform</InputLabel>
                            <Select name="platform" label="Platform" required>
                                {selectOptions.platforms.map((option: any, index) => {
                                    return (
                                        <MenuItem key={index} value={option.value}>{option.label}</MenuItem>
                                    )
                                })}
                            </Select>
                        </FormControl>

                        <DialogActions sx={{ mt: 1, mb: -1, mr: -1 }}>
                            <Button color="secondary" onClick={handleCreateGameDialogClose}>Cancel</Button>
                            <Button color="success" type="submit">Create</Button>
                        </DialogActions>
                    </form>
                </DialogContent>
            </Dialog>
        </Fragment>
    )
}

interface ExpandMoreProps extends IconButtonProps { expand: boolean }

const ExpandMore = styled((props: ExpandMoreProps) => {
    const { expand, ...other } = props;
    return <IconButton {...other} />;
})(({ theme, expand }) => ({
    transform: !expand ? 'rotate(0deg)' : 'rotate(180deg)',
    marginLeft: 'auto',
    transition: theme.transitions.create('transform', {
        duration: theme.transitions.duration.shortest,
    }),
}))